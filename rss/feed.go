// Package rss is to define generic model about rss feed
package rss

import "time"

type Feed struct {
	Title       string    `json:"title"`
	URL         string    `json:"url"`
	Source      string    `json:"source"`
	PublishedAt time.Time `json:"published_at"`
	Summary     string    `json:"summary"`
}
