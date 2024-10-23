package main

import (
	"errors"
	"fmt"
	"log"
)

func handlerLogin(st *state, cmd command) error {
	if len(cmd.Arguments) == 0 {
		return errors.New("login requires a username")
	}
	st.cfg.SetUser(cmd.Arguments[0])
	fmt.Printf("user has been set to %s.", cmd.Arguments[0])
	return nil
}

// register adds commands to the commands struct
func (c *commands) register(name string, f func(*state, command) error) {
	c.registeredCommands[name] = f
}

// run ensures a command has been registered, and then runs it
func (c *commands) run(s *state, cmd command) error {
	f, ok := c.registeredCommands[cmd.Name]
	if !ok {
		log.Fatalf("%s is not a command", cmd.Name)
	}

	return f(s, cmd)
}
