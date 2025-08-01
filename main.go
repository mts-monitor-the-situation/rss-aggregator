package main

import (
	"context"
	"fmt"
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

	rssURL := "https://moxie.foxnews.com/google-publisher/world.xml"
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
			Title:       item.Title,
			Description: item.Description,
			Link:        item.GetLink(),
			PubDate:     item.PubDate,
			Categories:  make([]mongodb.Category, 0, len(item.Categories)),
			GeoLocated:  false,
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
	err = feedItems.Save(rssCtx, collection)
	if err != nil {
		fmt.Println("error saving feed items:", err)
		return
	}
	fmt.Println("Feed items saved successfully")
}
