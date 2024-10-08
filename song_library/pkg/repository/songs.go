package repository

import (
	"github.com/ivnstd/SongLibrary/models"
	"github.com/sirupsen/logrus"
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
		logrus.Debugf("Applying filter: group = %s", group)
	}
	if song != "" {
		query = query.Where("song = ?", song)
		logrus.Debugf("Applying filter: song = %s", song)
	}
	if releaseDate != "" {
		query = query.Where("release_date = ?", releaseDate)
		logrus.Debugf("Applying filter: release_date = %s", releaseDate)
	}

	// Пагинация
	offset := (page - 1) * limit
	logrus.Debugf("Applying pagination: offset = %d, limit = %d", offset, limit)
	if err := query.Offset(offset).Limit(limit).Find(&songs).Error; err != nil {
		logrus.Errorf("Failed to execute query: %v", err)
		return nil, err
	}

	logrus.Debugf("Query executed successfully, retrieved %d songs", len(songs))
	return songs, nil
}

func (r *SongsDB) CreateSong(song models.Song) error {
	err := r.db.Create(&song).Error

	if err != nil {
		logrus.Errorf("Failed to execute query: %v", err)
		return err
	}

	logrus.Debug("Query executed successfully, song creates")
	return nil
}

func (r *SongsDB) GetSong(id uint) (models.Song, error) {
	var song models.Song
	err := r.db.First(&song, id).Error

	if err != nil {
		logrus.Errorf("Failed to execute query: %v", err)
		return song, err
	}

	logrus.Debugf("Query executed successfully, retrieved song with ID: %d", id)
	return song, err
}

func (r *SongsDB) UpdateSong(id uint, song models.Song) error {
	err := r.db.Model(&models.Song{}).Where("id = ?", id).Updates(song).Error

	if err != nil {
		logrus.Errorf("Failed to execute query: %v", err)
		return err
	}

	logrus.Debugf("Query executed successfully, updated song with ID: %d", id)
	return nil
}

func (r *SongsDB) DeleteSong(id uint) error {
	err := r.db.Delete(&models.Song{}, id).Error

	if err != nil {
		logrus.Errorf("Failed to execute query: %v", err)
		return err
	}

	logrus.Debugf("Query executed successfully, deleted song with ID: %d", id)
	return nil
}
