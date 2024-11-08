package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/ellielle/rssgator/internal/database"
)

// handlerAggregate fetches a feed and unescapes the feed string
func handlerAggregate(st *state, cmd command, user database.User) error {
	if len(cmd.Arguments) < 1 {
		return errors.New("usage: cli agg [time between requests](s/m/h)")
	}
	timeBetween, err := time.ParseDuration(cmd.Arguments[0])
	if err != nil {
		return errors.New("unable to parse duration, please use the format #s, #m, or #h")
	}
	fmt.Printf("collecting feeds every %v\n", timeBetween)
	ticker := time.NewTicker(timeBetween)
	for ; ; <-ticker.C {
		scrapeFeeds(st)
	}
}

// scrapeFeeds gets updated feed data for all stored feeds, at `duration` intervals
func scrapeFeeds(st *state) error {
	feedList, err := st.db.GetFeeds(context.Background())
	if err != nil {
		return errors.New("unable to retrieve feed list")
	}
	for range feedList {
		nextFeed, err := st.db.GetNextFeedFetch(context.Background())
		if err != nil {
			return errors.New("unable to get next feed from list")
		}
		feedData, err := fetchFeed(context.Background(), nextFeed.Url)
		if err != nil {
			return errors.New("unable to fetch feed")
		}
		err = st.db.MarkFeedFetched(context.Background(), nextFeed.ID)
		if err != nil {
			return errors.New("unable to mark feed fetched")
		}
		// unescape the feed text
		unescapeData(feedData)
		for _, item := range feedData.Channel.Item {
			publishedAt := sql.NullTime{}
			_, err := st.db.CreatePost(context.Background(), database.CreatePostParams{
				ID:        uuid.New(),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				Title:     item.Title,
				Url:       item.Link,
				Description: sql.NullString{
					String: item.Description,
					Valid:  true,
				},
				PublishedAt: publishedAt,
				FeedID:      nextFeed.ID,
			})
			if err != nil {
				if strings.Contains(err.Error(), "duplicate key value") {
					continue
				}
				log.Printf("couldn't create post: %v", err)
				continue
			}
		}
		log.Printf("feed %s collected, %v posts found", nextFeed.Name, len(feedData.Channel.Item))
	}
	return nil
}
