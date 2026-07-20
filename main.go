package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"

	"my-tech-radio/rss"
	"my-tech-radio/rss/zenn"
)

const candidatesFilePath = "candidates.json"

func main() {
	ctx := context.Background()
	feeds := make([]*rss.Feed, 0)

	client := http.DefaultClient
	zennAWSRSS, err := zenn.NewClient(client).FetchFeeds(ctx, "aws")
	if err != nil {
		slog.ErrorContext(ctx, "failed to fetech feeds from zenn", "error", err)
		panic(err)
	}

	zennAWSFeeds, err := zenn.ItemsToFeeds(zennAWSRSS.Channel.Items)
	if err != nil {
		slog.ErrorContext(ctx, "failed to convert items to feeds", "error", err)
		panic(err)
	}

	feeds = append(feeds, zennAWSFeeds...)

	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(feeds); err != nil {
		slog.ErrorContext(ctx, "failed to encoding feeds", "error", err)
		panic(err)
	}

	if err := os.WriteFile(candidatesFilePath, buf.Bytes(), 0644); err != nil {
		slog.ErrorContext(ctx, "failed to write candidates file", "error", err)
		panic(err)
	}
}
