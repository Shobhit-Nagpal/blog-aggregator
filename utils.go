package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"log"

	"github.com/Shobhit-Nagpal/blog-aggregator/internal/db"
	"github.com/Shobhit-Nagpal/blog-aggregator/internal/rss"
)

func printUsers(users []db.User, currentUser string) {
	for _, user := range users {
		if currentUser == user.Name {
			fmt.Printf("* %s (current)\n", user.Name)
			continue
		}

		fmt.Printf("* %s\n", user.Name)
	}
}

func scrapeFeeds(s *state) {
	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		log.Fatalf("Error getting next feed")
	}

	now := time.Now()
	feedParams := db.MarkFeedFetchedParams{
		ID:        feed.ID,
		UpdatedAt: now,
		LastFetchedAt: sql.NullTime{
			Time:  now,
			Valid: true,
		},
	}

	err = s.db.MarkFeedFetched(context.Background(), feedParams)
	if err != nil {
		log.Fatalf("Error marking feed as fetched")
	}

	rssFeed, err := rss.FetchFeed(context.Background(), feed.Url)
	if err != nil {
		log.Fatalf("Error fetching rss feed")
	}

	for idx, item := range rssFeed.Channel.Item {
		fmt.Printf("Item %d from %s: %s\n", idx+1, feed.Name, item.Title)
	}
}
