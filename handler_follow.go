package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/ellielle/rssgator/internal/database"
)

// handlerAddFollow adds a followed feed to a logged in user
func handlerAddFollow(st *state, cmd command, user database.User) error {
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
	// feed_follow isn't necessary here, _ it
	_, err = st.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return errors.New(err.Error())
	}
	fmt.Println("Feed follow created:")

	return nil
}

// handlerFollowing returns the names of all feeds the logged
// in user is following
func handlerFollowing(st *state, cmd command, user database.User) error {
	user, err := st.db.GetUserByName(context.Background(), st.cfg.CurrentUserName)
	if err != nil {
		return errors.New(err.Error())
	}
	ff, err := st.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return errors.New(err.Error())
	}
	for _, feed := range ff {
		f, err := st.db.GetFeedById(context.Background(), feed.FeedID)
		if err != nil {
			return errors.New(err.Error())
		}
		fmt.Println(f.Name)
	}

	return nil
}
