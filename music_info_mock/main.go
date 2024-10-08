package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type SongDetail struct {
	ReleaseDate string `json:"releaseDate"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}

func main() {
	r := gin.Default()

	r.GET("/info", func(c *gin.Context) {
		group := c.Query("group")
		song := c.Query("song")

		if group == "" || song == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "group and song are required"})
			return
		}

		response := SongDetail{
			ReleaseDate: "01.01.2010",
			Text:        "Ooh some text?\nOoh lalala?\n...\nSome text?\n\nOoh\nSome text\nOoh\nlalala",
			Link:        "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
		}

		c.JSON(http.StatusOK, response)
	})

	r.Run(":8081")
}
