package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Shobhit-Nagpal/blog-aggregator/internal/db"
	"github.com/google/uuid"
)

const FEED_URL = "https://www.wagslane.dev/index.xml"

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("No username provided")
	}

	user, err := s.db.GetUser(context.Background(), cmd.args[0])
	if err != nil {
		return err
	}

	s.cfg.SetUser(user.Name)

	fmt.Println("User has been set!")

	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("No name provided")
	}

	id := uuid.New()
	created_at := time.Now()
	updated_at := time.Now()

	params := db.CreateUserParams{
		ID:        id,
		CreatedAt: created_at,
		UpdatedAt: updated_at,
		Name:      cmd.args[0],
	}

	user, err := s.db.CreateUser(context.Background(), params)
	if err != nil {
		return err
	}

	s.cfg.SetUser(cmd.args[0])

	fmt.Println("User created successfully!")
	fmt.Println(user)

	return nil
}

func handlerReset(s *state, cmd command) error {
	err := s.db.DeleteAllUsers(context.Background())
	if err != nil {
		return err
	}

	fmt.Println("All users have been deleted")

	return nil
}

func handlerUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return err
	}

	printUsers(users, s.cfg.CurrentUserName)

	return nil
}

func handlerAggregate(s *state, cmd command) error {

	if len(cmd.args) == 0 {
		return errors.New("Missing args -> time between req")
	}

  timeBetweenRequests, err := time.ParseDuration(cmd.args[0])
  if err != nil {
    return err
  }

	ticker := time.NewTicker(timeBetweenRequests)

  fmt.Printf("Collecting feeds every %s\n\n", cmd.args[0])

	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
}

func handlerAddFeed(s *state, cmd command, user db.User) error {

	if len(cmd.args) < 2 {
		return errors.New("Missing args to add feed")
	}

	name := cmd.args[0]
	url := cmd.args[1]

	id := uuid.New()
	created_at := time.Now()
	updated_at := time.Now()

	params := db.CreateFeedParams{
		ID:        id,
		CreatedAt: created_at,
		UpdatedAt: updated_at,
		Name:      name,
		Url:       url,
		UserID:    user.ID,
	}

	feed, err := s.db.CreateFeed(context.Background(), params)
	if err != nil {
		return err
	}

	feedFollowParams := db.CreateFeedFollowParams{
		ID:        id,
		CreatedAt: created_at,
		UpdatedAt: updated_at,
		UserID:    user.ID,
		FeedID:    feed.ID,
	}
	_, err = s.db.CreateFeedFollow(context.Background(), feedFollowParams)
	if err != nil {
		return err
	}

	fmt.Println(feed)

	return nil
}

func handlerFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return err
	}

	for _, feed := range feeds {
		user, err := s.db.GetUserById(context.Background(), feed.UserID)
		if err != nil {
			return err
		}

		fmt.Printf("----Feed----\n\nName: %s\nURL: %s\nUsername: %s\n\n", feed.Name, feed.Url, user.Name)
	}

	return nil
}

func handlerFollow(s *state, cmd command, user db.User) error {

	if len(cmd.args) == 0 {
		return errors.New("Missing args to follow feed")
	}

	url := cmd.args[0]

	feed, err := s.db.GetFeedByUrl(context.Background(), url)
	if err != nil {
		return err
	}

	id := uuid.New()
	created_at := time.Now()
	updated_at := time.Now()

	params := db.CreateFeedFollowParams{
		ID:        id,
		CreatedAt: created_at,
		UpdatedAt: updated_at,
		UserID:    user.ID,
		FeedID:    feed.ID,
	}

	_, err = s.db.CreateFeedFollow(context.Background(), params)
	if err != nil {
		return err
	}

	fmt.Println(feed.Name)
	fmt.Println(user)

	return nil
}

func handlerFollowing(s *state, cmd command, user db.User) error {

	followFeed, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return err
	}

	for _, f := range followFeed {
		feed, err := s.db.GetFeedById(context.Background(), f.FeedID)
		if err != nil {
			return err
		}

		fmt.Println(feed.Name)
	}
	return nil
}

func handlerUnfollow(s *state, cmd command, user db.User) error {

	if len(cmd.args) == 0 {
		return errors.New("Missing args to unfollow feed")
	}

	url := cmd.args[0]

	feed, err := s.db.GetFeedByUrl(context.Background(), url)
	if err != nil {
		return err
	}

	params := db.DeleteFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	}

	err = s.db.DeleteFeedFollow(context.Background(), params)
	if err != nil {
		return err
	}

	return nil
}
