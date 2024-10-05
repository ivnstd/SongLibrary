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
			ReleaseDate: "16.07.2006",
			Text:        "Ooh baby, don't you know I suffer?\\nOoh baby, can you hear me moan?\\nYou caught me under false pretenses\\nHow long before you let me go?\\n\\nOoh\\nYou set my soul alight\\nOoh\\nYou set my soul alight",
			Link:        "https://www.youtube.com/watch?v=Xsp3_a-PMTw",
		}

		c.JSON(http.StatusOK, response)
	})

	r.Run(":8081")
}
