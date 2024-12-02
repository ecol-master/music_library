package update

import (
	"encoding/json"
	"log/slog"
	"music_lib/internal/utils"
	"net/http"
)

type SongUpdater interface {
	UpdateSong(utils.UpdatedSong) error
}

func New(songUpdater SongUpdater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.song.update"

		var req utils.UpdatedSong
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			slog.Error(op, "error decoding request", err)
			http.Error(w, "error decoding request", http.StatusBadRequest)
			return
		}

		err = songUpdater.UpdateSong(req)
		if err != nil {
			slog.Error(op, "error updating song", err)
			http.Error(w, "error updating song", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
