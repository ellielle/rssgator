package main

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"html"
	"io"
	"net/http"
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
