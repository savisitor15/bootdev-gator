package app

import (
	"context"
	"encoding/xml"
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

func fetchFeed(ctx context.Context, feedUrl string) (*RSSFeed, error) {
	var body io.Reader
	var client http.Client = http.Client{}
	req, err := http.NewRequestWithContext(ctx, "GET", feedUrl, body)
	if err != nil {
		return nil, fmt.Errorf("request failed to generate: %w", err)
	}
	req.Header.Set("User-Agent", "gator")
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode > 299 {
		return nil, fmt.Errorf("request failed: %d - %w", resp.StatusCode, err)
	}
	defer resp.Body.Close()
	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error parsing body to raw bytes: %w", err)
	}
	var rssOut RSSFeed
	err = xml.Unmarshal(raw, &rssOut)
	if err != nil {
		return nil, fmt.Errorf("error parsing from xml to feedbody: %w", err)
	}
	return &rssOut, nil
}

func cleanFeedOutput(inFeed *RSSFeed) error {
	// clean the header
	inFeed.Channel.Title = html.UnescapeString(inFeed.Channel.Title)
	inFeed.Channel.Description = html.UnescapeString(inFeed.Channel.Description)
	// clean the items
	for id, _ := range inFeed.Channel.Item {
		inFeed.Channel.Item[id].Title = html.UnescapeString(inFeed.Channel.Item[id].Title)
		inFeed.Channel.Item[id].Description = html.UnescapeString(inFeed.Channel.Item[id].Description)
	}
	return nil
}
