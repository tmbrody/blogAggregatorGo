package main

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/tmbrody/blogAggregatorGo/internal/database"
)

func processFeed(feedURL string, ctx context.Context, apiCfg *apiConfig) {
	rss, err := FetchRSSFeed(feedURL)
	if err != nil {
		fmt.Printf("Error fetching feed: %v", err)
		return
	}

	feeds, err := apiCfg.DB.GetAllFeeds(ctx)
	if err != nil {
		fmt.Printf("Error: %v", err)
		return
	}

	for _, feed := range feeds {
		if feed.Url == feedURL {
			apiCfg.Feed = feed
			break
		}
	}

	for _, item := range rss.Channel.Items {
		posts, err := apiCfg.DB.GetAllPosts(ctx)
		if err != nil {
			fmt.Printf("Error: %v", err)
			break
		}

		publishedAtString := item.PubDate

		var publishedAt sql.NullTime

		if publishedAtString != "" {
			var err error
			layouts := []string{
				time.RFC1123Z,
				"Mon, 02 Jan 2006 15:04:05 -0700",
				"Mon, 2 Jan 2006 15:04:05 -0700",
				"Mon, 02 Jan 2006 15:04:05 -0700 -0700",
				"Mon, 2 Jan 2006 15:04:05 -0700 -0700",
				"Mon, 02 Jan 2006 15:04:05 +0700 +0700",
				"Mon, 2 Jan 2006 15:04:05 +0700 +0700",
				"Mon, 02 Jan 2006 15:04:05 MST",
				"Mon, 2 Jan 2006 15:04:05 MST",
				"2006-01-02 15:04:05 -0700 -0700",
				"2006-01-02 15:04:05 +0700 +0700",
				"2006-01-02T15:04:05-07:00",
				"2006-01-02T15:04:05Z",
			}
			for _, layout := range layouts {
				publishedAt.Time, err = time.Parse(layout, publishedAtString)
				if err == nil {
					publishedAt.Valid = true
					break
				}
			}
			if !publishedAt.Valid {
				fmt.Printf("Error parsing time: %v", err)
			}
		} else {
			publishedAt.Valid = false
		}

		descriptionString := item.Description

		var description sql.NullString

		if descriptionString != "" {
			description.String = descriptionString
			description.Valid = true
		} else {
			description.Valid = false
		}

		postExists := false
		for _, post := range posts {
			if post.Title == item.Title.Value && post.Description == description {
				postExists = true
				break
			}
		}

		if postExists {
			continue
		}

		postID, err := uuid.NewUUID()
		if err != nil {
			fmt.Printf("Error: %v", err)
			break
		}

		title := strings.TrimSpace(item.Title.Value)

		args := database.CreatePostParams{
			ID:          postID,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       title,
			Url:         item.Link,
			Description: description,
			PublishedAt: publishedAt,
			FeedID:      apiCfg.Feed.ID,
		}

		_, err = apiCfg.DB.CreatePost(ctx, args)
		if err != nil {
			fmt.Printf("Error: %v", err)
			break
		}
	}
}

// Start a worker that will fetch and process feeds.
func feedWorker(ctx context.Context, feedChannel chan string, wg *sync.WaitGroup, apiCfg *apiConfig) {
	defer wg.Done()

	for {
		select {
		case feedURL := <-feedChannel:
			processFeed(feedURL, ctx, apiCfg)
		case <-ctx.Done():
			return
		}
	}
}
