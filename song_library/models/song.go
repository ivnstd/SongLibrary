package models

func (Song) TableName() string {
	return "song"
}

// @Description Песня в библиотеке
type Song struct {
	ID          uint   `json:"id"           gorm:"primaryKey;autoIncrement"` // ID песни
	Group       string `json:"group"        gorm:"not null"`                 // Исполнитель
	Song        string `json:"song"         gorm:"not null"`                 // Название песни
	ReleaseDate string `json:"releaseDate"`                                  // Дата релиза
	Text        string `json:"text"`                                         // Текст песни
	Link        string `json:"link"`                                         // Ссылка на песню
}

// @Description Данные, необходимые для добавления новой песни
type SongInput struct {
	Group string `json:"group" binding:"required"` // Исполнитель
	Song  string `json:"song"  binding:"required"` // Название песни
}

// @Description Подробная информация о песне, получаемая с внешнего API
type SongDetail struct {
	ReleaseDate string `json:"releaseDate"` // Дата релиза
	Text        string `json:"text"`        // Текст песни
	Link        string `json:"link"`        // Ссылка на песню
}
