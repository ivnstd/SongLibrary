package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/ivnstd/SongLibrary/pkg/service"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	songs := router.Group("/songs")
	{
		songs.GET("", h.get_songs)  // Получение данных библиотеки с фильтрацией по всем полям и пагинацией
		songs.POST("", h.post_song) // Добавление новой песни

		songsByID := songs.Group("/:id")
		{
			songsByID.Use(h.songMiddleware)

			songsByID.GET("", h.get_song)               // Получение полной информации о песне
			songsByID.PUT("", h.put_song)               // Изменение данных песни
			songsByID.DELETE("", h.delete_song)         // Удаление песни
			songsByID.GET("/lyrics", h.get_song_lyrics) // Получение текста песни с пагинацией по куплетам
		}
	}

	return router
}
