package repository

import (
	"github.com/ivnstd/SongLibrary/models"
	"gorm.io/gorm"
)

type Songs interface {
	GetSongs(group, song, releaseDate string, page, limit int) ([]models.Song, error)
	CreateSong(song models.Song) error
	GetSong(id uint) (models.Song, error)
	UpdateSong(id uint, song models.Song) error
	DeleteSong(id uint) error
}

type Repository struct {
	Songs
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		Songs: NewSongsDB(db),
	}
}
