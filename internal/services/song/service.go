package song

import (
	"music_lib/internal/api"
	"music_lib/internal/entities"
	"music_lib/internal/repositories/song"
	"music_lib/internal/utils"

	"github.com/pkg/errors"
)

type Service struct {
	repo *song.Repository
}

func New(repo *song.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) AddSong(group, song string, client *api.Client) (uint64, error) {
	const op = "song.Service.AddSong"

	songInfo, err := client.FetchSongStatic(group, song)
	if err != nil {
		return 0, errors.Wrap(err, op)
	}

	newSong := entities.Song{
		Group:       group,
		Song:        song,
		ReleaseDate: songInfo.ReleaseDate,
		Text:        songInfo.Text,
		Link:        songInfo.Link,
	}

	return s.repo.InsertSong(newSong)
}

func (s *Service) GetSong(id uint64) (*entities.Song, error) {
	const op = "song.Service.GetSong"
	return s.repo.GetSong(id)
}

func (s *Service) GetSongs(cursor_id, page_size uint64) ([]entities.Song, error) {
	const op = "song.Service.GetSongs"
	return s.repo.GetSongs(cursor_id, page_size)
}

func (s *Service) DeleteSong(id uint64) (uint64, error) {
	const op = "song.Service.DeleteSong"
	return s.repo.DeleteSong(id)
}

func (s *Service) UpdateSong(updatedSong utils.UpdatedSong) error {
	const op = "song.Service.UpdateSong"
	return s.repo.UpdateSong(updatedSong)
}

func (s *Service) FilterSongs(songFilter utils.FilteredSong) ([]entities.Song, error) {
	return s.repo.FilterSongs(songFilter)
}
