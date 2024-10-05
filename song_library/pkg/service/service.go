package service

import (
	"github.com/ivnstd/SongLibrary/models"
	"github.com/ivnstd/SongLibrary/pkg/repository"
)

type Songs interface {
	GetSongs(group, song, releaseDate string, page, limit int) ([]models.Song, error)
	CreateSong(song models.Song) error
	GetSong(id uint) (models.Song, error)
	UpdateSong(id uint, song models.Song) error
	DeleteSong(id uint) error
	GetSongLyrics(song models.Song, verseNumber int) (string, error)
}

type Service struct {
	Songs
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Songs: NewSongsService(repos.Songs),
	}
}
