package utils

import "time"

// All field except ID are optional
// If a field is nil, it will not be updated
type UpdatedSong struct {
	ID          uint       `json:"id" db:"id"`
	Group       *string    `json:"group" db:"group"`
	Song        *string    `json:"song" db:"song"`
	ReleaseDate *time.Time `json:"release_date" db:"release_date"`
	Text        *string    `json:"text" db:"text"`
	Link        *string    `json:"link" db:"link"`
}

type FilteredSong struct {
	Group       *string    `json:"group" db:"group"`
	Song        *string    `json:"song" db:"song"`
	ReleaseDate *time.Time `json:"release_date" db:"release_date"`
	Text        *string    `json:"text" db:"text"`
	Link        *string    `json:"link" db:"link"`
}
