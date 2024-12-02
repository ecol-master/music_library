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

type response struct {
	Song entities.Song `json:"song"`
}

type SongGetter interface {
	GetSong(uint64) (*entities.Song, error)
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

		err = json.NewEncoder(w).Encode(response{Song: *song})
		if err != nil {
			slog.Error(op, "error encoding response", err)
			http.Error(w, "error encoding response", http.StatusInternalServerError)
			return
		}
	}
}
