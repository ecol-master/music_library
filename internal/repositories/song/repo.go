package song

import (
	"github.com/jmoiron/sqlx"
	"music_lib/internal/entities"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) InsertSong(song entities.Song) (uint64, error) {
	var id uint64
	var q = `
		INSERT INTO "songs" ("group", "song", "release_date", "text", "link")
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`
	err := r.db.QueryRowx(q, song.Group, song.Song, song.ReleaseDate, song.Text, song.Link).Scan(&id)
	return id, err
}

func (r *Repository) DeleteSong(id uint64) (uint64, error) {
	var q = `
		DELETE FROM "songs"
		WHERE id = $1
		RETURNING id
	`
	err := r.db.QueryRowx(q, id).Scan(&id)
	return id, err
}
