package song

import (
	"github.com/jmoiron/sqlx"
	"music_lib/internal/entities"
)

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *repository {
	return &repository{db: db}
}

func (r *repository) InsertSong(song entities.Song) error {
	return nil
}
