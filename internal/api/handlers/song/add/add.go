package add

import (
	"encoding/json"
	"log/slog"
	"music_lib/internal/api"
	"net/http"
)

type request struct {
	Group string `json:"group"`
	Song  string `json:"song"`
}

type response struct {
	ID uint64 `json:"id"`
}

type SongAdder interface {
	AddSong(group, song string, client *api.Client) (uint64, error)
}

func New(songAdder SongAdder, client *api.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.song.add"

		var req request
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			slog.Error(op, "error decoding request", err)
			http.Error(w, "error decoding request", http.StatusBadRequest)
			return
		}

		id, err := songAdder.AddSong(req.Group, req.Song, client)
		if err != nil {
			slog.Error(op, "error adding song", err)
			http.Error(w, "error adding song", http.StatusInternalServerError)
			return
		}

		err = json.NewEncoder(w).Encode(response{ID: id})
		if err != nil {
			slog.Error(op, "error encoding response", err)
			http.Error(w, "error encoding response", http.StatusInternalServerError)
			return
		}
	}
}
