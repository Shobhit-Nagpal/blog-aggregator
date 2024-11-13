package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Shobhit-Nagpal/blog-aggregator/internal/db"
	"github.com/google/uuid"
)

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
