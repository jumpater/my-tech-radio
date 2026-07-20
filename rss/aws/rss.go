package aws

import (
	"encoding/xml"
	"fmt"
	"time"

	"my-tech-radio/rss"
)

type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Channel Channel  `xml:"channel"`
}

type Channel struct {
	Title string `xml:"title"`
	Items []Item `xml:"item"`
}

type Item struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PublishedAt string `xml:"pubDate"`
}

func ItemsToFeeds(items []Item) ([]rss.Feed, error) {
	feeds := make([]rss.Feed, 0, len(items))
	for _, item := range items {
		// AWSのpubDateは "Mon, 13 Jul 2026 18:13:57 +0000" 形式。
		// Zennの GMT 表記とは違い数値オフセットなのでRFC1123Zを使う。
		publishedAt, err := time.Parse(time.RFC1123Z, item.PublishedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to parse pubDate of %q: %w", item.Title, err)
		}

		feeds = append(feeds, rss.Feed{
			Title:       item.Title,
			URL:         item.Link,
			Source:      sourceName,
			PublishedAt: publishedAt,
			Summary:     item.Description,
		})
	}

	return feeds, nil
}
