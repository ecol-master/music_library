package song

import (
	"fmt"
	"music_lib/internal/entities"
	"music_lib/internal/utils"
	"strings"

	"github.com/jmoiron/sqlx"
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

func (r *Repository) GetSongs() ([]entities.Song, error) {
	var songs []entities.Song
	var q = `
		SELECT * FROM "songs"
	`
	err := r.db.Select(&songs, q)
	return songs, err
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

func (r *Repository) FilterSongs(songFilter utils.FilteredSong) ([]entities.Song, error) {
	var songs []entities.Song

	var queries []string
	var args []interface{}

	if songFilter.Group != nil {
		args = append(args, *songFilter.Group)
		queries = append(queries, fmt.Sprintf(`"group" = $%d`, len(args)))
	}

	if songFilter.Song != nil {
		args = append(args, *songFilter.Song)
		queries = append(queries, fmt.Sprintf(`"song" = $%d`, len(args)))
	}

	if songFilter.ReleaseDate != nil {
		args = append(args, *songFilter.ReleaseDate)
		queries = append(queries, fmt.Sprintf(`"release_date" = $%d`, len(args)))
	}

	if songFilter.Text != nil {
		args = append(args, *songFilter.Text)
		queries = append(queries, fmt.Sprintf(`"text" = $%d`, len(args)))
	}

	if songFilter.Link != nil {
		args = append(args, *songFilter.Link)
		queries = append(queries, fmt.Sprintf(`"link" = $%d`, len(args)))
	}

	if len(queries) == 0 {
		return nil, nil
	}

	query := fmt.Sprintf(`SELECT * FROM "songs" WHERE %s`, strings.Join(queries, " AND "))

	err := r.db.Select(&songs, query, args...)
	return songs, err
}
