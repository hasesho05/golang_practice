package controller

import (
	"go-sqlx-gin/db_client"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type History struct {
	User      User   `json:"user"`
	Word      string `json:"word"`
	CreatedAt string `json:"created_at"`
}

func CreateHistory(c *gin.Context) {
	var resBody struct {
		token string
		Word  string `json:"word"`
	}
	if err := c.ShouldBindJSON(&resBody); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "Invalid request body",
		})
		return
	}

	user, err := GetUserByToken(resBody.token)

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	res, err := db_client.DBClient.Exec("INSERT INTO history (user, word) VALUES (?, ?);",
		user.Id,
		resBody.Word,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	id, _ := res.LastInsertId()
	c.JSON(http.StatusCreated, gin.H{
		"error": false,
		"id":    id,
	})
}
