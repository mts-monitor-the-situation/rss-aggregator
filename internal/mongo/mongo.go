package mongo

import (
	"context"
	"fmt"
	"time"

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

type FeedItem struct {
	ID          string   `bson:"_id,omitempty"`
	Title       string   `bson:"title"`
	Description string   `bson:"description"`
	Link        string   `bson:"link"`
	PubDate     string   `bson:"pubDate"`
	Categories  []string `bson:"categories"`
	GeoLocated  bool     `bson:"geoLocated"`
	Latitude    float64  `bson:"latitude"`
	Longitude   float64  `bson:"longitude"`
}
