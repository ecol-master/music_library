package get

import (
	"encoding/json"
	"log/slog"
	"music_lib/internal/entities"
	"net/http"
)

type request struct {
	ID uint64 `json:"id"`
}

type SongGetter interface {
	GetSong(id uint64) (*entities.Song, error)
	GetSongs(cursor_id uint64, page_size uint64) ([]entities.Song, error)
}

func New(songUpdater SongGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.song.get"

		var req request
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			slog.Error(op, "error decoding request", err)
			http.Error(w, "error decoding request", http.StatusBadRequest)
			return
		}

		song, err := songUpdater.GetSong(req.ID)
		if err != nil {
			slog.Error(op, "error get song", err)
			http.Error(w, "error get song", http.StatusInternalServerError)
			return
		}

		err = json.NewEncoder(w).Encode(song)
		if err != nil {
			slog.Error(op, "error encoding response", err)
			http.Error(w, "error encoding response", http.StatusInternalServerError)
			return
		}
	}
}

type requestPaginated struct {
	CursorId uint64 `json:"cursor_id"`
	PageSize uint64 `json:"page_size"`
}

type responsePaginated struct {
	Songs    []entities.Song `json:"songs"`
	CursorId uint64          `json:"cursor_id"`
}

func NewAll(songsGetter SongGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.song.get_all"

		var req requestPaginated
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			slog.Error(op, "error decoding request", err)
			http.Error(w, "error decoding request", http.StatusBadRequest)
			return
		}

		songs, err := songsGetter.GetSongs(req.CursorId, req.PageSize)
		if err != nil {
			slog.Error(op, "error get song", err)
			http.Error(w, "error get song", http.StatusInternalServerError)
			return
		}

		cursor_id := req.CursorId
		if len(songs) > 0 {
			cursor_id = songs[len(songs)-1].ID + 1
		}

		resp := responsePaginated{
			Songs:    songs,
			CursorId: cursor_id,
		}

		err = json.NewEncoder(w).Encode(resp)
		if err != nil {
			slog.Error(op, "error encoding response", err)
			http.Error(w, "error encoding response", http.StatusInternalServerError)
			return
		}
	}
}
