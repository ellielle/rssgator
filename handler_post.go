package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/ellielle/rssgator/internal/database"
)

func handlerBrowsePosts(st *state, cmd command, user database.User) error {
	limit := 2
	if len(cmd.Arguments) > 0 {
		var err error
		limit, err = strconv.Atoi(cmd.Arguments[0])
		if err != nil {
			return err
		}
	}
	posts, err := st.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	})
	if err != nil {
		return fmt.Errorf("couldn't get posts for user: %w", err)
	}
	fmt.Printf("found %d posts for user %s:\n", len(posts), user.Name)
	for _, p := range posts {
		fmt.Printf("%s from %s\n", p.PublishedAt.Time.Format("Mon Jan 2"), p.FeedName)
		fmt.Printf("--- %s ---\n", p.Title)
		fmt.Printf("	%v\n", p.Description.String)
		fmt.Printf("Link: %s\n", p.Url)
		fmt.Println("====================================")
	}
	return nil
}
