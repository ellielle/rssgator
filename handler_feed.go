package main

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"html"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/ellielle/rssgator/internal/database"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func handlerAggregate(st *state, cmd command) error {
	//  if len(cmd.Arguments) == 0 {
	//  	fmt.Println("agg requires a URL as an argument")
	//  	os.Exit(1)
	//  }
	feed, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return errors.New("unable to fetch feed")
	}
	unescapeData(feed)
	fmt.Print(feed)
	return nil
}

func handlerAddFeed(st *state, cmd command) error {
	// Arguments[0] is the name of the feed, Arguments[1] is the URL of the feed
	if len(cmd.Arguments) < 2 {
		fmt.Println("addfeed requires a name and a URL")
	}
	user, err := st.db.GetUserByName(context.Background(), st.cfg.CurrentUserName)
	if err != nil {
		return err
	}
	_, err = st.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.Arguments[0],
		Url:       cmd.Arguments[1],
		UserID:    user.ID})
	if err != nil {
		return err
	}
	return nil
}

func handlerGetFeeds(st *state, cmd command) error {
	feeds, err := st.db.GetFeeds(context.Background())
	if err != nil {
		return err
	}
	for _, feed := range feeds {
		user, err := st.db.GetUserCreatedFeed(context.Background(), feed.UserID)
		if err != nil {
			return err
		}
		fmt.Println(feed.Name)
		fmt.Println(feed.Url)
		fmt.Println(user)
	}
	return nil
}

func unescapeData(feed *RSSFeed) error {
	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)
	for _, f := range feed.Channel.Item {
		f.Title = html.UnescapeString(f.Title)
		f.Description = html.UnescapeString(f.Description)
	}
	return nil
}

// fetchFeed gets an RSS feed and unmarshals it into a struct
func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return &RSSFeed{}, err
	}
	req.Header.Add("User-Agent", "gator")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return &RSSFeed{}, err
	}
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return &RSSFeed{}, err
	}
	feed := &RSSFeed{}
	xml.Unmarshal(data, &feed)
	return feed, nil
}
