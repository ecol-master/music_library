package entities

import "time"

type Song struct {
	ID          uint64       `json:"id" db:"id"`
	Group       string     `json:"group" db:"group"`
	Song        string     `json:"song" db:"song"`
	ReleaseDate *time.Time `json:"release_date" db:"release_date"`
	Text        string     `json:"text" db:"text"`
	Link        string     `json:"link" db:"link"`
}
