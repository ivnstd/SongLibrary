package models

func (Song) TableName() string {
	return "songs"
}

// @Description Песня в библиотеке
type Song struct {
	ID          uint   `json:"id"           gorm:"primaryKey;autoIncrement"` // ID песни
	Title       string `json:"title"        gorm:"not null"`                 // Название песни
	ArtistID    uint   `json:"-"            gorm:"not null"`                 // ID исполнителя
	ArtistName  string `json:"artist"`                                       // Наименование исполнителя
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
