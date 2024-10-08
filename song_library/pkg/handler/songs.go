package handler

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ivnstd/SongLibrary/models"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// Метод для получения данных библиотеки с фильтрацией по всем полям и пагинацией
// @Summary Get list of songs
// @Description Возвращает список песен с возможностью фильтрации по группе, названию песни и дате релиза, поддерживает пагинацию.
// @ID get-songs
// @Accept  json
// @Produce  json
// @Param group query string false "Фильтрация по группе или исполнителю"
// @Param song query string false "Фильтрация по названию песни"
// @Param release_date query string false "Фильтрация по дате релиза"
// @Param page query int false "Номер страницы для пагинации" default(1)
// @Param limit query int false "Количество элементов на странице" default(10)
// @Router /songs [get]
func (h *Handler) get_songs(c *gin.Context) {
	logrus.Infof("Handling get_songs request with params: group=%s, song=%s, release_date=%s, page=%s, limit=%s",
		c.DefaultQuery("group", ""), c.DefaultQuery("song", ""), c.DefaultQuery("release_date", ""), c.DefaultQuery("page", "1"), c.DefaultQuery("limit", "10"))

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
	logrus.Debug("Fetching songs")
	songs, err := h.services.Songs.GetSongs(group, title, releaseDate, page, limit)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve songs", err.Error())
		return
	}

	// Возврат данных о песнях
	newSuccessResponse(c, http.StatusOK, "Songs", songs)
}

// Метод для добавления новой песни
// @Summary Add new song
// @Description Добавление новой песни в библиотеку
// @ID post-song
// @Accept  json
// @Produce  json
// @Param song body models.SongInput true "Song Input"
// @Router /songs [post]
func (h *Handler) post_song(c *gin.Context) {
	logrus.Infof("Handling post_song request")

	var input models.SongInput

	// Проверка тела запроса на валидность
	if err := c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Invalid Song format")
		return
	}

	// Отправка запроса к внешнему API для получения дополнительной информации о песне
	logrus.Debugf("Fetching song details from external API with params: group %s, song: %s", input.Group, input.Song)
	songDetail, err := h.services.Songs.FetchSongDetail(input.Group, input.Song)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "Failed to fetch song details", err.Error())
		return
	}

	// Создание объекта песни с полученными данными
	song := models.Song{
		Group:       input.Group,
		Song:        input.Song,
		ReleaseDate: songDetail.ReleaseDate,
		Text:        songDetail.Text,
		Link:        songDetail.Link,
	}

	// Создание в базе записи, проверка на возможность создания
	logrus.Debugf("Creating song: %+v", song)
	if err := h.services.Songs.CreateSong(song); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "Failed to create song", err.Error())
		return
	}

	// Возврат сообщения о создании персоны
	newSuccessResponse(c, http.StatusCreated, "message", "Song created")
}

// Middleware для извлечения ID и проверки существования песни
func (h *Handler) songMiddleware(c *gin.Context) {
	idParam := c.Param("id")
	logrus.Debugf("Extracted ID from request: %s", idParam)

	// Проверка id на валидность
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Invalid ID format")
		c.Abort()
		return
	}

	// Проверка существования песни
	logrus.Debugf("Fetching song with ID: %d", id)
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
	logrus.Infof("Successfully retrieved song with ID: %d", id)
	c.Set("song", song)
	c.Next()
}

// Метод для получения данных песни
// @Summary Get song
// @Description Возвращает данные о песне по ее ID
// @ID get-song
// @Accept  json
// @Produce  json
// @Param id path int true "Song ID"
// @Router /songs/{id} [get]
func (h *Handler) get_song(c *gin.Context) {
	logrus.Infof("Handling get_song request")

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
// @Summary Update song
// @Description Изменяет данные о песни по её ID.
// @ID update-song
// @Accept  json
// @Produce  json
// @Param id path int true "Song ID"
// @Param song body models.Song true "Данные обновленной песни"
// @Router /songs/{id} [put]
func (h *Handler) put_song(c *gin.Context) {
	logrus.Infof("Handling put_song request")

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
	logrus.Debugf("Updating song with ID: %d, %+v", song.ID, updatedSong)
	if err := h.services.Songs.UpdateSong(song.ID, updatedSong); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "Failed to update person", err.Error())
		return
	}

	// Возврат сообщения об обновлении песни
	newSuccessResponse(c, http.StatusOK, "message", "Song updated")
}

// Метод для удаления песни
// @Summary Delete song
// @Description Удаляет существующую песню по её ID.
// @ID delete-song
// @Accept  json
// @Produce  json
// @Param id path int true "Song ID"
// @Router /songs/{id} [delete]
func (h *Handler) delete_song(c *gin.Context) {
	logrus.Infof("Handling delete_song request")

	// Получение данных о песне из контекста
	song := c.MustGet("song").(models.Song)

	// Удаление из базы записи, проверка на возможность удаления
	logrus.Debugf("Deleting song with ID: %d", song.ID)
	if err := h.services.Songs.DeleteSong(song.ID); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "Failed to delete song", err.Error())
		return
	}

	// Возврат сообщения об удалении песни
	newSuccessResponse(c, http.StatusOK, "message", "Song deleted")
}

// Метод для получения текста песни с пагинацией по куплетам
// @Summary Get song lyrics
// @Description Возвращает текст куплета песни по номеру куплета.
// @ID get-song-lyrics
// @Accept  json
// @Produce  json
// @Param id path int true "Song ID"
// @Param verse query int false "Номер куплета (по умолчанию 1)"
// @Router /songs/{id}/lyrics [get]
func (h *Handler) get_song_lyrics(c *gin.Context) {
	logrus.Infof("Handling get_song_lyrics request")

	// Получение данных о песне из контекста
	song := c.MustGet("song").(models.Song)

	// Извлечение параметра для пагинации
	verseStr := c.DefaultQuery("verse", "1")
	verseNumber, err := strconv.Atoi(verseStr)
	if err != nil || verseNumber < 1 {
		newErrorResponse(c, http.StatusBadRequest, "Invalid verse number")
		return
	}
	logrus.Debugf("Extracted verse number: %d", verseNumber)

	// Получение данных о куплете песни
	logrus.Debugf("Fetching song's lyrics with ID: %d", song.ID)
	verse, err := h.services.Songs.GetSongLyrics(song, verseNumber)
	if err != nil {
		newErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	// Возврат куплета песни
	newSuccessResponse(c, http.StatusOK, "Verse", verse)
}
