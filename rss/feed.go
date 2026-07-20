// Package rss は収集した記事の共通モデルを定義する。
// 取得とXMLの解析は取得元ごとのサブパッケージ(rss/zenn など)が担当する。
package rss

import "time"

// Feed は収集した記事1件。candidates.json の要素になる。
type Feed struct {
	Title       string    `json:"title"`
	URL         string    `json:"url"`
	Source      string    `json:"source"`
	PublishedAt time.Time `json:"published_at"`
	Summary     string    `json:"summary"`
}
