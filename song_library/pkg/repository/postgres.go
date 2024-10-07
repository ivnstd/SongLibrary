package repository

import (
	"fmt"

	"github.com/ivnstd/SongLibrary/models"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

func NewDB(cfg Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password, cfg.SSLMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&models.Song{})

	return db, nil
}

func SeedDatabaseIfEmpty(db *gorm.DB) {
	var count int64
	db.Model(&models.Song{}).Count(&count)
	if count == 0 {
		SeedDatabase(db)
	}
}

func SeedDatabase(db *gorm.DB) {
	songs := []models.Song{
		{
			Group:       "Muse",
			Song:        "Supermassive Black Hole",
			ReleaseDate: "16.07.2006",
			Text:        "Ooh baby, don't you know I suffer?\nOoh baby, can you hear me moan?\nYou caught me under false pretenses\nHow long before you let me go?\n\nOoh\nYou set my soul alight\nOoh\nYou set my soul alight...",
			Link:        "https://www.youtube.com/watch?v=Xsp3_a-PMTw",
		}, {
			Group:       "Muse",
			Song:        "Uprising",
			ReleaseDate: "04.08.2009",
			Text:        "Paranoia is in bloom, the PR transmissions will resume\nThey'll try to push drugs that keep us all dumbed down\nAnd hope that we will never see the truth around\n(So come on)\n\nAnother promise, another seed\nAnother packaged lie to keep us trapped in greed...",
			Link:        "https://www.youtube.com/watch?v=w8KQmps-Sog",
		}, {
			Group:       "Muse",
			Song:        "Hysteria",
			ReleaseDate: "01.12.2003",
			Text:        "It's bugging me, grating me\nAnd twisting me around\nYeah, I'm endlessly caving in\nAnd turning inside out\n\nCause I want it now, I want it now\nGive me your heart and your soul...\n\nYeah, it's holding me, morphing me...",
			Link:        "https://www.youtube.com/watch?v=3dm_5qWWDV8",
		}, {
			Group:       "twenty one pilots",
			Song:        "Tear in My Heart",
			ReleaseDate: "06.04.2015",
			Text:        "An-nyŏng-ha-se-yo\n\nSometimes you gotta bleed to know\nThat you're alive and have a soul\nBut it takes someone to come around to show you how\n\nShe's the tear in my heart, I'm alive\nShe's the tear in my heart, I'm on fire...\n\nThe songs on the radio are okay...",
			Link:        "https://www.youtube.com/watch?v=nky4me4NP70",
		}, {
			Group:       "twenty one pilots",
			Song:        "Heavydirtysoul",
			ReleaseDate: "17.05.2015",
			Text:        "There's an infestation in my mind's imagination\nI hope that they choke on smoke...\n\nGangsters don't cry, therefore, therefore I'm (I'm)\nMr. Misty-eyed, therefore I'm (I'm)...",
			Link:        "https://www.youtube.com/watch?v=r_9Kf0D5BTs",
		}, {
			Group:       "no name artist",
			Song:        "Heavydirtysoul",
			ReleaseDate: "14.02.2020",
			Text:        "some text...\n\nlalal\nsome text...",
			Link:        "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
		}, {
			Group:       "some artist",
			Song:        "my track",
			ReleaseDate: "16.07.2006",
			Text:        "some text...\n\nlalal\nsome text...",
			Link:        "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
		},
	}

	if err := db.Create(&songs).Error; err != nil {
		logrus.Fatalf("Failed to seed db: %s", err.Error())
	} else {
		logrus.Println("Database seeded with initial data")
	}
}
