package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/ivnstd/SongLibrary/models"
	"github.com/ivnstd/SongLibrary/pkg/repository"
)

type SongsService struct {
	repo repository.Songs
}

func NewSongsService(repo repository.Songs) *SongsService {
	return &SongsService{repo: repo}
}

func (s *SongsService) GetSongs(group, song, releaseDate string, page, limit int) ([]models.Song, error) {
	return s.repo.GetSongs(group, song, releaseDate, page, limit)
}

func (s *SongsService) FetchSongDetail(group string, song string) (*models.SongDetail, error) {
	url := fmt.Sprintf("http://localhost:8081/info?group=%s&song=%s", group, song)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Unexpected status code: %d", resp.StatusCode)
	}

	var songDetail models.SongDetail
	if err := json.NewDecoder(resp.Body).Decode(&songDetail); err != nil {
		return nil, fmt.Errorf("Failed to decode response: %v", err)
	}

	return &songDetail, nil
}

func (s *SongsService) CreateSong(song models.Song) error {
	return s.repo.CreateSong(song)
}

func (s *SongsService) GetSong(id uint) (models.Song, error) {
	return s.repo.GetSong(id)
}

func (s *SongsService) UpdateSong(id uint, song models.Song) error {
	return s.repo.UpdateSong(id, song)
}

func (s *SongsService) DeleteSong(id uint) error {
	return s.repo.DeleteSong(id)
}

func (s *SongsService) GetSongLyrics(song models.Song, verseNumber int) (string, error) {
	verses := strings.Split(song.Text, "\\n\\n")

	if verseNumber < 1 || verseNumber > len(verses) {
		return "", fmt.Errorf("Verse not found")
	}

	return verses[verseNumber-1], nil
}
