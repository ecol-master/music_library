package filter

import (
	"encoding/json"
	"log/slog"
	"music_lib/internal/entities"
	"music_lib/internal/utils"
	"net/http"
)

type SongsFilter interface {
	FilterSongs(utils.FilteredSong) ([]entities.Song, error)
}

func New(songFilter SongsFilter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.song.filter"

		var req utils.FilteredSong
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			slog.Error(op, "error decoding request", err)
			http.Error(w, "error decoding request", http.StatusBadRequest)
			return
		}

		songs, err := songFilter.FilterSongs(req)
		if err != nil {
			slog.Error(op, "error filtering song", err)
			http.Error(w, "error filtering song", http.StatusInternalServerError)
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
