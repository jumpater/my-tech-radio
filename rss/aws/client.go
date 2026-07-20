// Package aws はAWS News Blogのフィードから記事を収集する。
package aws

import (
	"context"
	"encoding/xml"
	"fmt"
	"net/http"

	"my-tech-radio/rss"
)

const (
	sourceName = "aws"
	feedURL    = "https://aws.amazon.com/blogs/aws/feed/"
)

func NewClient(client *http.Client) *Client {
	return &Client{client}
}

type Client struct {
	*http.Client
}

// FetchFeeds はAWS News Blogのフィードを取得して記事一覧を返す。
func (c *Client) FetchFeeds(ctx context.Context) ([]rss.Feed, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, feedURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch feed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %s", resp.Status)
	}

	var r RSS
	if err := xml.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, fmt.Errorf("failed to decode feed: %w", err)
	}

	return ItemsToFeeds(r.Channel.Items)
}
