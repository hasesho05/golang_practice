package controller

import (
	"database/sql"
	"go-sqlx-gin/db_client"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type Post struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Content   *string   `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

func CreatePost(c *gin.Context) {
	var resBody Post
	if err := c.ShouldBindJSON(&resBody); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "Invalid request body",
		})
		return
	}

	res, err := db_client.DBClient.Exec("INSERT INTO posts (title, content) VALUES (?, ?);",
		resBody.Title,
		resBody.Content,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": true,
		})
		return
	}

	id, _ := res.LastInsertId()
	c.JSON(http.StatusCreated, gin.H{
		"error": false,
		"id":    id,
	})
}

func GetPosts(c *gin.Context) {
	var posts []Post

	rows, err := db_client.DBClient.Query("SELECT id, title, content, created_at FROM posts;")
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "Invalid request body",
		})
		return
	}

	for rows.Next() {
		var singlePost Post
		if err := rows.Scan(&singlePost.ID, &singlePost.Title, &singlePost.Content, &singlePost.CreatedAt); err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"error":   true,
				"message": "Invalid request body",
			})
			return
		}
		posts = append(posts, singlePost)
	}

	c.JSON(http.StatusOK, posts)
}

func GetPost(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)
	row := db_client.DBClient.QueryRow("SELECT id, title, content, created_at FROM posts WHERE id = ?;", id)
	var post Post
	if err := row.Scan(&post.ID, &post.Title, &post.Content, &post.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"error": true,
			})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": true,
			})
		}
		return
	}
	c.JSON(http.StatusOK, post)
}
