package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/ellielle/rssgator/internal/database"
)

func handlerAddFollow(st *state, cmd command) error {
	if len(cmd.Arguments) < 1 {
		fmt.Println("following requires a url")
	}
	user, err := st.db.GetUserByName(context.Background(), st.cfg.CurrentUserName)
	if err != nil {
		return errors.New(err.Error())
	}
	feed, err := st.db.GetFeedByUrl(context.Background(), cmd.Arguments[0])
	if err != nil {
		return errors.New(err.Error())
	}
	ff, err := st.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	fmt.Printf("%s\n%s", feed.Name, st.cfg.CurrentUserName)

	return nil
}
