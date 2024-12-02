package song

import (
	"github.com/pkg/errors"
	"music_lib/internal/api"
	"music_lib/internal/entities"
	"music_lib/internal/repositories/song"
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
