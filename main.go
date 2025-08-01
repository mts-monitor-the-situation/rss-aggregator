package main

import (
	"context"
	"fmt"
	"html"
	"time"

	"github.com/mts-monitor-the-situation/rss-aggregator/internal/mongodb"
	"github.com/mts-monitor-the-situation/rss-aggregator/internal/redisdb"
	"github.com/mts-monitor-the-situation/rss-aggregator/pkg/rss"
)

func main() {

	// Connect to Redis
	redisClient, err := redisdb.Connect("localhost:6379")
	if err != nil {
		fmt.Printf("error connecting to Redis: %v", err)
		return
	}
	defer redisClient.Close()

	// Connect to MongoDB
	mongoURI := "mongodb://localhost:27017"
	client, err := mongodb.Connect(mongoURI)
	if err != nil {
		fmt.Printf("error connecting to MongoDB: %v", err)
		return
	}
	defer client.Disconnect(context.Background())

	// RSS logic
	rssCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// rssURL := "https://moxie.foxnews.com/google-publisher/world.xml"
	rssURL := "https://feeds.bbci.co.uk/news/world/rss.xml"
	// rssURL := "https://rss.nytimes.com/services/xml/rss/nyt/World.xml"
	rss, err := rss.FetchRSS(rssURL, rssCtx)
	if err != nil {
		fmt.Println("error fetching RSS:", err)
		return
	}

	// Print the RSS title
	fmt.Println("RSS Title:", rss.Channel.Title)

	// Create a feed items collection
	collection := client.Database("mts").Collection("feed_items")
	feedItems := mongodb.FeedItems{}

	// Process each item in the RSS feed
	for _, item := range rss.Channel.Items {
		feedItem := mongodb.FeedItem{
			ID:          item.GenId(),
			Source:      rss.Channel.Title,
			Title:       html.UnescapeString(item.Title),
			Description: html.UnescapeString(item.Description),
			Link:        item.GetLink(),
			PubDate:     item.PubDate,
			Categories:  make([]mongodb.Category, 0, len(item.Categories)),
		}

		// Convert categories
		for _, c := range item.Categories {
			feedItem.Categories = append(feedItem.Categories, mongodb.Category{
				Domain: c.Domain,
				Value:  c.Value,
			})
		}

		feedItems.Items = append(feedItems.Items, feedItem)
	}

	// Save feed items to MongoDB
	ids, err := feedItems.Save(rssCtx, collection)
	if err != nil {
		fmt.Println("error saving feed items:", err)
		return
	}

	unprocessedFeedLength := len(ids)
	if unprocessedFeedLength == 0 {
		fmt.Println("no feed items to save")
		return
	}

	fmt.Printf("saved %d feed items to MongoDB\n", unprocessedFeedLength)

	// Publish only upserted IDs to Redis stream
	for _, id := range ids {
		fmt.Printf("publishing ID to Redis stream: %s\n", id)
		err := redisdb.AddToStream(redisClient, "rss:unprocessed", map[string]any{
			"id": id,
		})
		if err != nil {
			fmt.Printf("failed to publish to Redis stream: %v", err)
		}
	}

	// // Now publish IDs to Redis for geolocation
	// for _, feedItem := range feedItems.Items {
	// 	if !feedItem.GeoLocated {
	// 		err := redisdb.AddToStream(redisClient, "rss:unprocessed", map[string]any{
	// 			"id": feedItem.ID,
	// 		})
	// 		if err != nil {
	// 			fmt.Printf("failed to publish to Redis stream: %v", err)
	// 		}
	// 	}
	// }

	fmt.Println("unprocessed feed items published to Redis stream successfully")
}
