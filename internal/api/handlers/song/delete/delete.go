package delete

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type Request struct {
	ID uint64 `json:"id"`
}

type Response struct {
	ID uint64 `json:"id"`
}

type SongDeleter interface {
	DeleteSong(id uint64) (uint64, error)
}

func New(songDeleter SongDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.song.delete"

		var req Request
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			slog.Error(op, "error decoding request", err)
			http.Error(w, "error decoding request", http.StatusBadRequest)
			return
		}

		id, err := songDeleter.DeleteSong(req.ID)
		if err != nil {
			slog.Error(op, "error deleting song", err)
			http.Error(w, "error deleting song", http.StatusInternalServerError)
			return
		}

		err = json.NewEncoder(w).Encode(Response{ID: id})
		if err != nil {
			slog.Error(op, "error encoding response", err)
			http.Error(w, "error encoding response", http.StatusInternalServerError)
			return
		}
	}
}
