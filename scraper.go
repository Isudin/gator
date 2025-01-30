package main

import (
	"context"
	"fmt"

	"github.com/Isudin/gator/feed"
	"github.com/google/uuid"
)

func scrapeFeeds(s *state) error {
	nextFeed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return err
	}

	if nextFeed.ID == uuid.Nil {
		return fmt.Errorf("no feeds found in database")
	}

	err = s.db.MarkFeedFetched(context.Background(), nextFeed.ID)
	if err != nil {
		return nil
	}

	fetchedFeed, err := feed.FetchFeed(context.Background(), nextFeed.Url)
	if err != nil {
		return err
	}

	fmt.Printf("--- %v\n:", fetchedFeed.Channel.Title)
	for _, item := range fetchedFeed.Channel.Item {
		fmt.Printf("%v, ", item.Title)
	}
	fmt.Println()

	return nil
}
