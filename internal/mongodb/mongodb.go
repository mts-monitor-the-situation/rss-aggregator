package mongodb

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// Connect establishes a connection to the MongoDB database
// Remember to defer client.Disconnect(ctx) after using the client
func Connect(uri string) (*mongo.Client, error) {

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)

	// Client options
	opts := options.Client()
	opts.ApplyURI(uri).SetServerAPIOptions(serverAPI)
	opts.SetConnectTimeout(10 * time.Second)

	client, err := mongo.Connect(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Ping the server to verify connection
	ctxPing, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = client.Ping(ctxPing, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	return client, nil
}

// GenId generates a deterministic ID based on guid (if present) and link
func GenId(guid string, link string, pubDate string) string {

	input := ""

	if strings.TrimSpace(guid) != "" {
		input = guid + link + pubDate
	} else {
		input = link + pubDate
	}

	hash := sha1.Sum([]byte(input)) // returns [20]byte

	return hex.EncodeToString(hash[:])
}

// FeedItem represents a single RSS feed item stored in MongoDB
type FeedItem struct {
	ID          string     `bson:"_id"`
	Source      string     `bson:"source"`
	Title       string     `bson:"title"`
	Description string     `bson:"description"`
	Link        string     `bson:"link"`
	PubDate     string     `bson:"pubDate"`
	Categories  []Category `bson:"categories"`
	GeoLocated  bool       `bson:"geoLocated"`
	Latitude    float64    `bson:"latitude"`
	Longitude   float64    `bson:"longitude"`
}

// Category represents a category in an RSS item, stored in MongoDB
type Category struct {
	Domain string `bson:"domain"`
	Value  string `bson:"value"`
}

// FeedItems is a wrapper for a slice of FeedItem for MongoDB operations
type FeedItems struct {
	Items []FeedItem `bson:"items"`
}

// Save saves the FeedItems to the MongoDB collection
func (f *FeedItems) Save(ctx context.Context, collection *mongo.Collection) error {
	if len(f.Items) == 0 {
		return nil
	}

	var models []mongo.WriteModel
	for _, item := range f.Items {
		model := mongo.NewUpdateOneModel().
			SetFilter(bson.M{"_id": item.ID}).
			SetUpdate(bson.M{"$set": item}).
			SetUpsert(true)

		models = append(models, model)
	}

	_, err := collection.BulkWrite(ctx, models, options.BulkWrite().SetOrdered(false))
	if err != nil {
		return fmt.Errorf("bulk write failed: %w", err)
	}
	return nil
}
