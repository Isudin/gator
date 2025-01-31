package main

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/Isudin/gator/feed"
	"github.com/Isudin/gator/internal/database"
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

	fmt.Printf("--- %v\n", fetchedFeed.Channel.Title)
	for _, item := range fetchedFeed.Channel.Item {
		savePosts(item, nextFeed, s)
	}
	fmt.Println()

	return nil
}

func savePosts(item feed.RSSItem, nextFeed database.Feed, s *state) {
	sqlTime := sql.NullTime{}
	parsedTime, err := time.Parse(time.RFC1123Z, item.PubDate)
	if err != nil {
		fmt.Printf("error parsing date '%v' - %v\n", item.PubDate, err)
	} else {
		sqlTime.Time = parsedTime
		sqlTime.Valid = true
	}

	sqlDescription := sql.NullString{}
	if item.Description != "" {
		sqlDescription.String = item.Description
		sqlDescription.Valid = true
	}

	params := database.CreatePostParams{
		ID:          uuid.New(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Title:       item.Title,
		Url:         item.Link,
		Description: sqlDescription,
		PublishedAt: sqlTime,
		FeedID:      nextFeed.ID,
	}
	_, err = s.db.CreatePost(context.Background(), params)
	if err != nil && !strings.HasPrefix(err.Error(), "pq: duplicate key value") {
		fmt.Println(err)
	}
}
