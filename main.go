package main

import (
	"database/sql"
	"errors"
	"log"
	"os"

	"github.com/Shobhit-Nagpal/blog-aggregator/internal/config"
	"github.com/Shobhit-Nagpal/blog-aggregator/internal/db"
	_ "github.com/lib/pq"
)

func main() {
	cfg := config.Read()

	database, err := sql.Open("postgres", cfg.DBUrl)
	if err != nil {
		log.Fatalf("Couldn't connect to database: %s\n", err.Error())
	}

	dbQueries := db.New(database)

	s := &state{
		db:  dbQueries,
		cfg: &cfg,
	}

	cmds := commands{
		handlers: map[string]func(*state, command) error{},
	}

	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
  cmds.register("reset", handlerReset)
  cmds.register("users", handlerUsers)
  cmds.register("agg", handlerAggregate)
  cmds.register("addfeed", handleAddFeed)

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
