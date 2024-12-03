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


// @Summary UpdateSong
// @Description Update song by ID
// @Tags update
// @Accept json
// @Produce json
// @Param request body utils.UpdatedSong true "Request"
// @Success 200 {string} string
// @Failure 400 {string} string
// @Failure 500 {string} string
// @Router /update_song [put]
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
