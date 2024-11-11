package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/Shobhit-Nagpal/blog-aggregator/internal/config"
)

func main() {
	cfg := config.Read()

	s := &state{
		cfg: &cfg,
	}

	cmds := commands{
		handlers: map[string]func(*state, command) error{},
	}

	cmds.register("login", handlerLogin)

	args := os.Args
	if len(args) < 2 {
		log.Fatalf("No command given")
	}

	cmd := command{
		name: os.Args[1],
		args: os.Args[2:],
	}

  err := cmds.run(s, cmd)
  if err != nil {
    log.Fatal(err)
  }

}

type state struct {
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

	s.cfg.SetUser(cmd.args[0])

	fmt.Println("User has been set!")

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
