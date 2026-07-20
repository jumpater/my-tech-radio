// Package zenn is feature for recieving zenn's rss feed
package zenn

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

func ItemsToFeeds(items []Item) ([]*rss.Feed, error) {
	feeds := make([]*rss.Feed, 0)
	for _, item := range items {
		publishedAt, err := parsePubDate(item.PublishedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to parse pubDate of %q: %w", item.Title, err)
		}

		feed := &rss.Feed{
			Title:       item.Title,
			URL:         item.Link,
			Source:      "zenn",
			PublishedAt: publishedAt,
			Summary:     item.Description,
		}

		feeds = append(feeds, feed)
	}

	return feeds, nil
}

// pubDateのゾーンは略称(GMT)と数値オフセット(+0900)の両方がありうるが、
// time.RFC1123は前者、time.RFC1123Zは後者しか受け付けない。
func parsePubDate(s string) (time.Time, error) {
	for _, layout := range []string{time.RFC1123, time.RFC1123Z} {
		if t, err := time.Parse(layout, s); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("unrecognized pubDate format %q", s)
}
