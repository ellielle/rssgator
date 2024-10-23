package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"

	"github.com/ellielle/rssgator/internal/database"
)

func handlerLogin(st *state, cmd command) error {
	if len(cmd.Arguments) == 0 {
		return errors.New("login requires a username")
	}

	// search for the user in the database, return an error if they don't exist
	user, err := st.db.GetUserByName(context.Background(), cmd.Arguments[0])
	if err != nil {
		fmt.Println("user does not exist")
		os.Exit(1)
	}

	st.cfg.SetUser(user.Name)
	fmt.Printf("user has been set to %s.\n", user.Name)
	return nil
}

func handlerRegister(st *state, cmd command) error {
	if len(cmd.Arguments) == 0 {
		return errors.New("register requires a username")
	}
	user, err := st.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		Name:      cmd.Arguments[0],
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		fmt.Println("error: user already exists")
		os.Exit(1)
	}
	// set current user in config to newly created user
	err = st.cfg.SetUser(user.Name)
	if err != nil {
		return errors.New("error saving new user to config")
	}
	fmt.Printf("user %s created!\n", user.Name)
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
