package song

import (
	"github.com/jmoiron/sqlx"
	"music_lib/internal/entities"
	"music_lib/internal/utils"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetSong(id uint64) (*entities.Song, error) {
	var song entities.Song
	var q = `
		SELECT * FROM "songs"
		WHERE id = $1
	`
	err := r.db.Get(&song, q, id)
	return &song, err
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

func (r *Repository) UpdateSong(updatedSong utils.UpdatedSong) error {
	originalSong, err := r.GetSong(uint64(updatedSong.ID))
	if err != nil {
		return err
	}

	if updatedSong.Group == nil {
		updatedSong.Group = &originalSong.Group
	}

	if updatedSong.Song == nil {
		updatedSong.Song = &originalSong.Song
	}

	if updatedSong.ReleaseDate == nil {
		updatedSong.ReleaseDate = originalSong.ReleaseDate
	}

	if updatedSong.Text == nil {
		updatedSong.Text = &originalSong.Text
	}

	if updatedSong.Link == nil {
		updatedSong.Link = &originalSong.Link
	}

	var q = `
		UPDATE "songs"
		SET "group" = $1, "song" = $2, "release_date" = $3, "text" = $4, "link" = $5
		WHERE id = $6
	`

	err = r.db.QueryRowx(q, updatedSong.Group, updatedSong.Song, updatedSong.ReleaseDate, updatedSong.Text, updatedSong.Link, updatedSong.ID).Err()
	return err
}
