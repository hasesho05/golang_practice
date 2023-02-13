package controller

import (
	"fmt"
	"go-sqlx-gin/db_client"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type History struct {
	User      int    `json:"user"`
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

	fmt.Println("userid", user.Id)

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

func GetHistory(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "ng",
			"message": "user not found",
		})
	}
	user, err := GetUserByToken(token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "ng",
			"message": err.Error(),
		})
		return
	}

	row, err := db_client.DBClient.Query("SELECT * FROM history WHERE user = ?", user.Id)
	if err != nil {
		c.JSON(http.StatusNoContent, gin.H{
			"status":  "ng",
			"message": err.Error(),
		})
		return
	}
	var history History

	err = row.Scan(&history.Word, &history.User, &history.CreatedAt)
	if err != nil {
		c.JSON(http.StatusNoContent, gin.H{
			"status":  "ng",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"status":  "ok",
		"message": "ログインしました。",
		"data":    history,
	})
}
