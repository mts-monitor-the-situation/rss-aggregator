package main

import (
	"context"
	"fmt"
	"time"

	"github.com/mts-monitor-the-situation/rss-aggregator/pkg/rss"
)

func main() {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Example usage of the RSS package
	rssURL := "https://moxie.foxnews.com/google-publisher/world.xml"
	rss, err := rss.FetchRSS(rssURL, ctx)
	if err != nil {
		fmt.Println("Error fetching RSS:", err)
		return
	}

	fmt.Println("RSS Title:", rss.Channel.Title)
	for _, item := range rss.Channel.Items {
		fmt.Printf("Item: %s, Link: %s\n", item.Title, item.GetLink())
	}
}
