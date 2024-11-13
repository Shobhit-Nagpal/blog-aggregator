package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Shobhit-Nagpal/blog-aggregator/internal/config"
	"github.com/Shobhit-Nagpal/blog-aggregator/internal/db"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

func main() {
	cfg := config.Read()

	conn, err := sql.Open("postgres", cfg.DBUrl)
	if err != nil {
		log.Fatalf("Couldn't connect to database: %s\n", err.Error())
	}

	dbQueries := db.New(conn)

	s := &state{
		db:  dbQueries,
		cfg: &cfg,
	}

	cmds := commands{
		handlers: map[string]func(*state, command) error{},
	}

	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)

	args := os.Args
	if len(args) < 2 {
		log.Fatalf("No command given")
	}

	cmd := command{
		name: os.Args[1],
		args: os.Args[2:],
	}

	err = cmds.run(s, cmd)
	if err != nil {
		log.Fatal(err)
	}

}

type state struct {
	db  *db.Queries
	cfg *config.Config
}

type command struct {
	name string
	args []string
}

type commands struct {
	handlers map[string]func(*state, command) error
}

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

func (c *commands) register(name string, f func(*state, command) error) {
	c.handlers[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	if handler, exists := c.handlers[cmd.name]; exists {
		err := handler(s, cmd)
		if err != nil {
			return err
		}
		return nil
	}

	return errors.New("Handler does not exist for commmand given")
}
