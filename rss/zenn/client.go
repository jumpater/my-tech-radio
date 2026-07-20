package zenn

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

func NewClient(client *http.Client) *Client {
	return &Client{
		client,
	}
}

type Client struct {
	*http.Client
}

func (c *Client) FetchFeeds(ctx context.Context, topic string) (*RSS, error) {
	var buf []byte
	reader := bytes.NewReader(buf)

	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("https://zenn.dev/topics/%s/feed", topic),
		reader,
	)
	if err != nil {
		slog.ErrorContext(ctx, "failed to create request", "error", err)
		return nil, err
	}

	resp, err := c.Do(req)
	if err != nil {
		slog.ErrorContext(ctx, "failed to fetch zenn rss feeds.", "error", err)
		return nil, err
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.ErrorContext(ctx, "failed to read response", "error", err)
		return nil, err
	}

	var rss RSS
	err = xml.Unmarshal(data, &rss)
	if err != nil {
		slog.ErrorContext(ctx, "failed to unmarshal xml", "error", err)
		return nil, err
	}

	if err := resp.Body.Close(); err != nil {
		return nil, err
	}

	return &rss, nil
}
