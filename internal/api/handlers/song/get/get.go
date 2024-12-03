package get

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"log/slog"
	"music_lib/internal/entities"
	"net/http"
	"strconv"

	_ "music_lib/docs"
)

type SongGetter interface {
	GetSong(id uint64) (*entities.Song, error)
	GetSongs(cursor_id uint64, page_size uint64) ([]entities.Song, error)

	GetSongText(id, cursor_id, offset uint64) ([]string, error)
}

// @Summary GetSong
// @Description Get song by ID
// @Tags get
// @Accept json
// @Produce json
// @Param id path uint64 true "Song ID" 
// @Success 200 {object} entities.Song
// @Failure 400 {string} string
// @Failure 500 {string} string
// @Router /get_song/{id} [get]  // путь без параметра в URL

func New(songUpdater SongGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.song.get"

		idParam := chi.URLParam(r, "id")
		slog.Debug(op, "id", idParam)
		id, err := strconv.ParseUint(idParam, 10, 64)

		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		song, err := songUpdater.GetSong(id)
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

// @Summary GetAllSongs
// @Description Get all songs with pagination
// @Tags get
// @Accept json
// @Produce json
// @Param request body requestPaginated true "Pagination parameters"
// @Success 200 {object} responsePaginated
// @Failure 400 {string} string
// @Failure 500 {string} string
// @Router /get_all_songs [get]
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

type requestGetPaginatedText struct {
	ID       uint64 `json:"id"`
	CursorId uint64 `json:"cursor_id"`
	Offset   uint64 `json:"offset"`
}

type responseGetPaginatedText struct {
	Text     []string `json:"text"`
	CursorId uint64   `json:"cursor_id"`
}

// @Summary GetSongText
// @Description Get song text by ID with pagination
// @Tags get
// @Accept json
// @Produce json
// @Param request body requestGetPaginatedText true "Pagination parameters"
// @Success 200 {object} responseGetPaginatedText
// @Failure 400 {string} string
// @Failure 500 {string} string
// @Router /get_song_text [get]
func NewText(songsGetter SongGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.song.get_all"

		var req requestGetPaginatedText
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			slog.Error(op, "error decoding request", err)
			http.Error(w, "error decoding request", http.StatusBadRequest)
			return
		}

		text, err := songsGetter.GetSongText(req.ID, req.CursorId, req.Offset)
		if err != nil {
			slog.Error(op, "error get song", err)
			http.Error(w, "error get song", http.StatusInternalServerError)
			return
		}

		resp := responseGetPaginatedText{
			Text:     text,
			CursorId: req.CursorId + req.Offset + 1,
		}

		err = json.NewEncoder(w).Encode(resp)
		if err != nil {
			slog.Error(op, "error encoding response", err)
			http.Error(w, "error encoding response", http.StatusInternalServerError)
			return
		}
	}
}
