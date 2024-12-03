package delete

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"log/slog"
	"net/http"
	"strconv"
)

type request struct {
	ID uint64 `json:"id"`
}

type response struct {
	ID uint64 `json:"id"`
}

type SongDeleter interface {
	DeleteSong(id uint64) (uint64, error)
}

// @Summary DeleteSong
// @Description Delete song by ID
// @Tags delete
// @Accept json
// @Produce json
// @Param request body request true "Request"
// @Success 200 {object} response
// @Failure 400 {string} string
// @Failure 500 {string} string
// @Router /delete_song/{id} [delete]
func New(songDeleter SongDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.song.delete_song"

		idParam := chi.URLParam(r, "id")
		slog.Debug(op, "id", idParam)
		id, err := strconv.ParseUint(idParam, 10, 64)

		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		id, err = songDeleter.DeleteSong(id)
		if err != nil {
			slog.Error(op, "error deleting song", err)
			http.Error(w, "error deleting song", http.StatusInternalServerError)
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
