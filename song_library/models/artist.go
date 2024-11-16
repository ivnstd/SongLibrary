package models

func (Artist) TableName() string {
	return "artists"
}

// @Description Исполнитель песен
type Artist struct {
	ID   uint   `json:"id"           gorm:"primaryKey;autoIncrement"` // ID исполнителя
	Name string `json:"name"       gorm:"not null"`                   // Наименование исполнителя

	Songs []Song `json:"-" gorm:"foreignKey:ArtistID;constraint:OnDelete:CASCADE;"`
}
