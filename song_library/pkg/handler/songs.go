package handler

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ivnstd/SongLibrary/models"
	"gorm.io/gorm"
)

// Метод для получения данных библиотеки с фильтрацией по всем полям и пагинацией
func (h *Handler) get_songs(c *gin.Context) {
	// Извлечение параметров для фильтрации и пагинации
	group := c.DefaultQuery("group", "")
	title := c.DefaultQuery("song", "")
	releaseDate := c.DefaultQuery("release_date", "")

	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}

	// Получение данных о песнях из базы
	songs, err := h.services.Songs.GetSongs(group, title, releaseDate, page, limit)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve songs", err.Error())
		return
	}

	// Возврат данных о песнях
	newSuccessResponse(c, http.StatusOK, "Songs", songs)
}

// Метод для создания новой песни
func (h *Handler) post_song(c *gin.Context) {
	var input models.SongInput

	// Проверка тела запроса на валидность
	if err := c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Invalid Song format")
		return
	}

	//TODO: запрос на внешний API/API MOCK, обработка ошибок
	song := models.Song{
		Group:       input.Group,
		Song:        input.Song,
		ReleaseDate: "16.07.2006",
		Text:        "Ooh baby, don't you know I suffer?\\nOoh baby, can you hear me moan?\\nYou caught me under false pretenses\\nHow long before you let me go?\\n\\nOoh\\nYou set my soul alight\\nOoh\\nYou set my soul alight",
		Link:        "https://www.youtube.com/watch?v=Xsp3_a-PMTw",
	}

	// Создание в базе записи, проверка на возможность создания
	if err := h.services.Songs.CreateSong(song); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "Failed to create song", err.Error())
		return
	}

	// Возврат сообщения о создании персоны
	newSuccessResponse(c, http.StatusCreated, "message", "Song created")
}

// Middleware для извлечения ID и проверки существования персоны
func (h *Handler) songMiddleware(c *gin.Context) {
	idParam := c.Param("id")

	// Проверка id на валидность
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Invalid ID format")
		c.Abort()
		return
	}

	// Проверка существования песни
	song, err := h.services.Songs.GetSong(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newErrorResponse(c, http.StatusNotFound, "Song not found")
		} else {
			newErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve song", err.Error())
		}
		c.Abort()
		return
	}

	// Сохранение данных о песне в контексте
	c.Set("song", song)
	c.Next()
}

// Метод для получения данных песни
func (h *Handler) get_song(c *gin.Context) {
	// Получение данных о песне из контекста
	song, exists := c.Get("song")
	if !exists {
		newErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve song from context")
		return
	}

	// Возврат данных о песне
	newSuccessResponse(c, http.StatusOK, "Song", song)
}

// Метод для изменения данных песни
func (h *Handler) put_song(c *gin.Context) {
	var updatedSong models.Song

	// Получение данных о песне из контекста
	song := c.MustGet("song").(models.Song)

	// Проверка тела запроса на валидность
	if err := c.ShouldBindJSON(&updatedSong); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Invalid Song format")
		return
	}

	// Дополнительная проверка формата даты
	_, err := time.Parse("02.01.2006", updatedSong.ReleaseDate)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Invalid Song format")
		return
	}

	// Обновление данных в базе, проверка на возможность обновления
	if err := h.services.Songs.UpdateSong(song.ID, updatedSong); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "Failed to update person", err.Error())
		return
	}

	// Возврат сообщения об обновлении песни
	newSuccessResponse(c, http.StatusOK, "message", "Song updated")
}

// Метод для удаления песни
func (h *Handler) delete_song(c *gin.Context) {
	// Получение данных о песне из контекста
	song := c.MustGet("song").(models.Song)

	// Удаление из базы записи, проверка на возможность удаления
	if err := h.services.Songs.DeleteSong(song.ID); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "Failed to delete song", err.Error())
		return
	}

	// Возврат сообщения об удалении песни
	newSuccessResponse(c, http.StatusOK, "message", "Song deleted")
}

// Метод для получения текста песни с пагинацией по куплетам
func (h *Handler) get_song_lyrics(c *gin.Context) {
	// Получение данных о песне из контекста
	song := c.MustGet("song").(models.Song)

	// ПИзвлечение параметра для пагинации
	verseStr := c.DefaultQuery("verse", "1")
	verseNumber, err := strconv.Atoi(verseStr)
	if err != nil || verseNumber < 1 {
		newErrorResponse(c, http.StatusBadRequest, "Invalid verse number")
		return
	}

	// Получение данных о куплете песни
	verse, err := h.services.Songs.GetSongLyrics(song, verseNumber)
	if err != nil {
		newErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	// Возврат куплета песни
	newSuccessResponse(c, http.StatusOK, "Verse", verse)
}
