package models

import "time"

type Item struct {
	Title       string
	Link        string
	Description string
	PubDate     time.Time
}
