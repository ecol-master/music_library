package filter

import (
	"encoding/json"
	"log/slog"
	"music_lib/internal/entities"
	"music_lib/internal/utils"
	"net/http"
)

type requestPaginated struct {
	Filters  utils.FilteredSong `json:"filters"`
	CursorId uint64             `json:"cursor_id"`
	PageSize uint64             `json:"page_size"`
}

type responsePaginated struct {
	Songs    []entities.Song `json:"songs"`
	CursorId uint64          `json:"cursor_id"`
}

type SongsFilter interface {
	FilterSongs(filteredSong utils.FilteredSong, cursor_id uint64, page_size uint64) ([]entities.Song, error)
}

// @Summary FilterSongs
// @Description Filter songs by Song fields
// @Tags filter
// @Accept json
// @Produce json
// @Param request body requestPaginated true "Request"
// @Success 200 {object} responsePaginated
// @Failure 400 {string} string
// @Failure 500 {string} string
// @Router /filter_songs [get]
func New(songFilter SongsFilter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.song.filter"

		var req requestPaginated
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			slog.Error(op, "error decoding request", err)
			http.Error(w, "error decoding request", http.StatusBadRequest)
			return
		}

		songs, err := songFilter.FilterSongs(req.Filters, req.CursorId, req.PageSize)

		if err != nil {
			slog.Error(op, "error filtering song", err)
			http.Error(w, "error filtering song", http.StatusInternalServerError)
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
