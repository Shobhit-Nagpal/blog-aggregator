package main

import (
	"context"

	"github.com/Shobhit-Nagpal/blog-aggregator/internal/db"
)

func middlewareLoggedIn(handler func(s *state, cmd command, user db.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
		if err != nil {
			return err
		}

		err = handler(s, cmd, user)
		if err != nil {
			return err
		}

		return nil
	}
}
