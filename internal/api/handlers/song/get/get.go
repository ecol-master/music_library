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
	GetSong(uint64) (*entities.Song, error)
	GetSongs() ([]entities.Song, error)
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

func NewAll(songsGetter SongGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.song.get_all"

		songs, err := songsGetter.GetSongs()
		if err != nil {
			slog.Error(op, "error get song", err)
			http.Error(w, "error get song", http.StatusInternalServerError)
			return
		}

		err = json.NewEncoder(w).Encode(songs)
		if err != nil {
			slog.Error(op, "error encoding response", err)
			http.Error(w, "error encoding response", http.StatusInternalServerError)
			return
		}
	}
}
