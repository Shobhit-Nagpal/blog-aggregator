package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"log"

	"github.com/Shobhit-Nagpal/blog-aggregator/internal/db"
	"github.com/Shobhit-Nagpal/blog-aggregator/internal/rss"
	"github.com/google/uuid"
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
    fmt.Printf("Item %d from %s: %s ; Published: %s\n", idx+1, feed.Name, item.Title, item.PubDate)

		id := uuid.New()
		created_at := time.Now()
		updated_at := time.Now()
		published_at, err := time.Parse(time.RFC1123, item.PubDate)
		if err != nil {
			log.Fatalf("Error parsing published date")
		}

		params := db.CreatePostParams{
			ID:        id,
			CreatedAt: created_at,
			UpdatedAt: updated_at,
			Title:     item.Title,
			Url:       item.Link,
			Description: sql.NullString{
				String: item.Description,
				Valid:  true,
			},
			PublishedAt: published_at,
      FeedID: feed.ID,
		}

		err = s.db.CreatePost(context.Background(), params)
		if err != nil {
      if err.Error() != `pq: duplicate key value violates unique constraint "posts_url_key"` {
      log.Fatalf("Error creating post: %s\n", err.Error())
      }
		}
	}
}
