package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/tmbrody/blogAggregatorGo/internal/database"
)

func scrapeWorker(dbQueries *database.Queries) {
	sleepDuration := 1 * time.Minute

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	feedChannel := make(chan string)
	fmt.Println("Starting feed worker")

	var wg sync.WaitGroup
	wg.Add(1)
	go feedWorker(ctx, feedChannel, &wg)

	for {
		feeds, err := dbQueries.GetNextFeedsToFetch(ctx)
		if err != nil {
			fmt.Printf("Error: %v", err)
			break
		}

		for _, feed := range feeds {
			feedChannel <- feed.Url

			_, err := dbQueries.MarkFeedFetched(ctx, feed.ID)
			if err != nil {
				fmt.Printf("Error: %v", err)
				break
			}
		}

		time.Sleep(3 * time.Second)

		fmt.Println("The next set of feeds will be fetched after 60 seconds...")
		fmt.Println()

		time.Sleep(sleepDuration)
	}

	close(feedChannel)
	wg.Wait()
}
