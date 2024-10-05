package models

func (Song) TableName() string {
	return "song"
}

type Song struct {
	ID          uint   `json:"id"           gorm:"primaryKey;autoIncrement"`
	Group       string `json:"group"        gorm:"not null"`
	Song        string `json:"song"         gorm:"not null"`
	ReleaseDate string `json:"releaseDate"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}

type SongInput struct {
	Group string `json:"group" binding:"required"`
	Song  string `json:"song"  binding:"required"`
}

type SongDetail struct {
	ReleaseDate string `json:"releaseDate"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}
