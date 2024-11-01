package main

import (
	"context"
	"errors"

	"github.com/ellielle/rssgator/internal/database"
)

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		user, err := s.db.GetUserByName(context.Background(), s.cfg.CurrentUserName)
		if err != nil {
			return errors.New(err.Error())
		}

		return handler(s, cmd, user)

	}

}
