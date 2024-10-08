package main

import (
	db_server "github.com/ivnstd/SongLibrary"
	"github.com/ivnstd/SongLibrary/configs"
	"github.com/ivnstd/SongLibrary/pkg/handler"
	"github.com/ivnstd/SongLibrary/pkg/repository"
	"github.com/ivnstd/SongLibrary/pkg/service"
	"github.com/sirupsen/logrus"

	_ "github.com/ivnstd/SongLibrary/docs"
)

// @title Song Library API
// @version 1.0
// @description Реализация онлайн библиотеки песен.
// @host localhost:8080
// @BasePath /

func main() {
	// logrus.SetLevel(logrus.InfoLevel)
	logrus.SetLevel(logrus.DebugLevel)

	if err := configs.LoadConfig(); err != nil {
		logrus.Fatalf("Error loading env variables: %s", err.Error())
	}

	logrus.Info("Starting server...")

	db, err := repository.NewDB(repository.Config{
		Host:     configs.Config.DB_Host,
		Port:     configs.Config.DB_Port,
		Username: configs.Config.DB_Username,
		DBName:   configs.Config.DB_Name,
		SSLMode:  configs.Config.DB_SSLMode,
		Password: configs.Config.DB_Password,
	})
	if err != nil {
		logrus.Fatalf("Failed to initialize db: %s", err.Error())
	}
	logrus.Info("Database connection established")

	repository.SeedDatabaseIfEmpty(db)

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	router := handlers.InitRoutes()

	srv := new(db_server.Server)
	if err := srv.Run(configs.Config.MainPort, router); err != nil {
		logrus.Fatalf("Error occured while running http server: %s", err.Error())
	}
	logrus.Info("http server successfully launched")
}
