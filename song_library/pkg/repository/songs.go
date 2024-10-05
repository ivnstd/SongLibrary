package repository

import (
	"github.com/ivnstd/SongLibrary/models"
	"gorm.io/gorm"
)

type SongsDB struct {
	db *gorm.DB
}

func NewSongsDB(db *gorm.DB) *SongsDB {
	return &SongsDB{db: db}
}

func (r *SongsDB) GetSongs(group, song, releaseDate string, page, limit int) ([]models.Song, error) {
	var songs []models.Song
	query := r.db.Model(&models.Song{})

	// Фильтрация по параметрам
	if group != "" {
		query = query.Where("\"group\" = ?", group)
	}
	if song != "" {
		query = query.Where("song = ?", song)
	}
	if releaseDate != "" {
		query = query.Where("release_date = ?", releaseDate)
	}

	// Пагинация
	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Find(&songs).Error; err != nil {
		return nil, err
	}

	return songs, nil
}

func (r *SongsDB) CreateSong(song models.Song) error {
	err := r.db.Create(&song).Error
	return err
}

func (r *SongsDB) GetSong(id uint) (models.Song, error) {
	var song models.Song
	err := r.db.First(&song, id).Error
	return song, err
}

func (r *SongsDB) UpdateSong(id uint, song models.Song) error {
	err := r.db.Model(&models.Song{}).Where("id = ?", id).Updates(song).Error
	return err
}

func (r *SongsDB) DeleteSong(id uint) error {
	err := r.db.Delete(&models.Song{}, id).Error
	return err
}
