package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"slices"

	"my-tech-radio/rss"
	"my-tech-radio/rss/aws"
	"my-tech-radio/rss/zenn"
)

const candidatesFilePath = "candidates.json"

// 収集対象のZennトピック。増やしたいときはここに足す。
var zennTopics = []string{"aws", "mysql", "ai", "go"}

func main() {
	ctx := context.Background()

	if err := run(ctx); err != nil {
		slog.ErrorContext(ctx, "failed to collect feeds", "error", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	client := http.DefaultClient
	feeds := make([]rss.Feed, 0)

	zennClient := zenn.NewClient(client)
	for _, topic := range zennTopics {
		fetched, err := zennClient.FetchFeeds(ctx, topic)
		if err != nil {
			// 1つ落ちても残りは集めたいので、警告を出して次に進む
			slog.WarnContext(ctx, "skipped zenn topic", "topic", topic, "error", err)
			continue
		}

		feeds = append(feeds, fetched...)
	}

	awsFeeds, err := aws.NewClient(client).FetchFeeds(ctx)
	if err != nil {
		slog.WarnContext(ctx, "skipped aws blog", "error", err)
	} else {
		feeds = append(feeds, awsFeeds...)
	}

	// 新しい記事から並べる
	slices.SortFunc(feeds, func(a, b rss.Feed) int {
		return b.PublishedAt.Compare(a.PublishedAt)
	})

	data, err := json.MarshalIndent(feeds, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to encode feeds: %w", err)
	}

	if err := os.WriteFile(candidatesFilePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write %s: %w", candidatesFilePath, err)
	}

	slog.InfoContext(ctx, "wrote candidates", "path", candidatesFilePath, "count", len(feeds))

	return nil
}
