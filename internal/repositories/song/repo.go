package song

import (
	"fmt"
	"log/slog"
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

func (r *Repository) GetSongs(cursor_id, page_size uint64) ([]entities.Song, error) {
	var songs []entities.Song
	var q = `
        SELECT * FROM songs
        WHERE id >= $1
        ORDER BY id ASC
        LIMIT $2
    `
	err := r.db.Select(&songs, q, cursor_id, page_size)
	return songs, err
}

func (r *Repository) GetSongText(id, cursor_id, offset uint64) ([]string, error) {
	var text string
	var q = `
		SELECT "text" FROM "songs"
		WHERE id = $1
	`
	err := r.db.Get(&text, q, id)
	slog.Debug("Song text", "text", text)

	if err != nil {
		return nil, err
	}

	data := strings.Split(text, "\n\n")
	slog.Debug("Song text", "data", data, "len", len(data), "cursor_id", cursor_id, "offset", offset)

	if len(data) <= int(cursor_id) {
		return []string{}, nil
	}

	var textArr []string
	if len(data) > int(cursor_id+offset) {
		textArr = data[cursor_id : cursor_id+offset]
	} else {
		textArr = data[cursor_id:]
	}
	return textArr, nil
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

func (r *Repository) FilterSongs(songFilter utils.FilteredSong, cursor_id, page_size uint64) ([]entities.Song, error) {
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

	query := fmt.Sprintf(`
		SELECT * FROM "songs" 
		WHERE %s AND id >= $%d
        ORDER BY id ASC
        LIMIT $%d
		`, strings.Join(queries, " AND "), len(args)+1, len(args)+2)

	args = append(args, cursor_id, page_size)
	err := r.db.Select(&songs, query, args...)
	return songs, err
}
