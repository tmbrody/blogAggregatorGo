package main

import (
	"context"
	"fmt"
	"strings"
	"sync"
)

// Fetch and process a feed (for documentation, print the post titles).
func processFeed(feedURL string) {
	rss, err := FetchRSSFeed(feedURL)
	if err != nil {
		fmt.Printf("Error fetching feed: %v", err)
	}

	fmt.Printf("Processing feed: %v\n", rss.Channel.Title)
	// Print the titles of the posts in the feed's channel
	for _, item := range rss.Channel.Items {
		title := strings.TrimSpace(item.Title.Value)
		fmt.Printf("\tTitle: %v\n", title)
	}
	fmt.Println("------------------------------------------------------------")
}

// Start a worker that will fetch and process feeds.
func feedWorker(ctx context.Context, feedChannel chan string, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case feedURL := <-feedChannel:
			processFeed(feedURL)
		case <-ctx.Done():
			return
		}
	}
}
