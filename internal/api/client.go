package api

import (
	"encoding/json"
	"fmt"
	"music_lib/internal/config"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

type SongInfo struct {
	ReleaseDate *time.Time `json:"release_date" db:"release_date"`
	Text        string     `json:"text" db:"text"`
	Link        string     `json:"link" db:"link"`
}

type Client struct {
	HTTPClient *http.Client
	config     config.APIClientConfig
}

func NewClient(clientConfig config.APIClientConfig) *Client {
	return &Client{
		HTTPClient: &http.Client{
			Timeout: clientConfig.RequestTimeout,
		},
		config: clientConfig,
	}
}

func (c *Client) FetchSong(group, song string) (*SongInfo, error) {
	url := fmt.Sprintf("%s?group=%s&song=%s", c.config.BaseURL, group, song)
	response, err := c.HTTPClient.Get(url)

	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	switch response.StatusCode {
	case http.StatusOK:
		var songInfo SongInfo
		err = json.NewDecoder(response.Body).Decode(&songInfo)
		if err != nil {
			return nil, errors.Wrap(err, "failed to decode response body")
		}
		return &songInfo, nil

	// Handle failed cases
	case http.StatusBadRequest:
		return nil, errors.New("bad request")
	case http.StatusInternalServerError:
		return nil, errors.New("internal server error")
	default:
		return nil, errors.New("unexpected status code")
	}
}

// Only needed for testing purposes
func (c *Client) FetchSongStatic(group, song string) (*SongInfo, error) {
	return &SongInfo{
		ReleaseDate: new(time.Time),
		Text:        "text",
		Link:        "link",
	}, nil
}
