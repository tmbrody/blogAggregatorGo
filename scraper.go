package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/tmbrody/blogAggregatorGo/internal/database"
)

func (apiCfg *apiConfig) scrapeWorker(dbQueries *database.Queries) {
	sleepDuration := 1 * time.Minute

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	feedChannel := make(chan string)

	var wg sync.WaitGroup
	wg.Add(1)
	go feedWorker(ctx, feedChannel, &wg, apiCfg)

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

		time.Sleep(sleepDuration)
	}

	close(feedChannel)
	wg.Wait()
}
