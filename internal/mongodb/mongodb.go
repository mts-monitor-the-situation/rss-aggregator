package mongodb

import (
	"context"
	"fmt"
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

	// Create a new client and connect to the server
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
	Locations   []Location `bson:"locations"`
}

// Category represents a category in an RSS item, stored in MongoDB
type Category struct {
	Domain string `bson:"domain"`
	Value  string `bson:"value"`
}

// Location represents a geographical location of an RSS item, stored in MongoDB
type Location struct {
	Latitude  float64 `bson:"latitude"`
	Longitude float64 `bson:"longitude"`
	PlaceID   string  `bson:"placeId"`
}

// FeedItems is a wrapper for a slice of FeedItem for MongoDB operations
type FeedItems struct {
	Items []FeedItem `bson:"items"`
}

// Save saves the FeedItems to the MongoDB collection
func (f *FeedItems) Save(ctx context.Context, collection *mongo.Collection) ([]string, error) {

	// Ensure there are items to save
	if len(f.Items) == 0 {
		return nil, nil
	}

	// Prepare bulk write models for upsert
	models := make([]mongo.WriteModel, 0, len(f.Items))

	for _, item := range f.Items {
		model := mongo.NewUpdateOneModel().
			SetFilter(bson.M{"_id": item.ID}).
			SetUpdate(bson.M{"$set": item}).
			SetUpsert(true)

		models = append(models, model)
	}

	// Perform the bulk write operation
	bkRes, err := collection.BulkWrite(ctx, models, options.BulkWrite().SetOrdered(false))
	if err != nil {
		return nil, fmt.Errorf("bulk write failed: %w", err)
	}

	// Collect the IDs of upserted items
	ids := make([]string, 0, len(bkRes.UpsertedIDs))
	for _, upsertedID := range bkRes.UpsertedIDs {
		if idStr, ok := upsertedID.(string); ok {
			ids = append(ids, idStr)
		}
	}

	return ids, nil
}
