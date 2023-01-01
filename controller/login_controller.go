package controller

import (
	"go-sqlx-gin/db_client"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"
)

type User struct {
	id         *int64  `json:id`
	Username   string  `json:"username"`
	Password   string  `json:"password"`
	Prefecture *string `json:"prefecture"`
}

// 暗号(Hash)化
func PasswordEncrypt(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

// 暗号(Hash)と入力された平パスワードの比較
func CompareHashAndPassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func CreateUser(c *gin.Context) {
	var resBody User
	if err := c.ShouldBindJSON(&resBody); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "Invalid request body",
			"content": err,
		})
		return
	}

	encryptedPassword, err := PasswordEncrypt(resBody.Password)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "Invalid request body",
			"content": err.Error(),
		})
	}

	res, err := db_client.DBClient.Exec("INSERT INTO user (username, password, prefecture) VALUES (?, ?, ?);",
		resBody.Username,
		encryptedPassword,
		resBody.Prefecture,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"content": err.Error(),
			"message": "既に存在するユーザーネームです。",
		})
		return
	}

	id, _ := res.LastInsertId()
	c.JSON(http.StatusCreated, gin.H{
		"status": "ok",
		"error":  false,
		"id":     id,
	})
}

func Login(c *gin.Context) {
	var resBody User
	if err := c.ShouldBindJSON(&resBody); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "Invalid request body",
			"content": err,
		})
		return
	}
	dbPassword := GetUser(resBody.Username).Password
	formPassword := resBody.Password

	if err := CompareHashAndPassword(dbPassword, formPassword); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "ng",
			"message": "パスワードが違います。",
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{
		"status":  "ok",
		"message": "ログインしました。",
	})
}

func GetUser(username string) User {
	row := db_client.DBClient.QueryRow("SELECT id, username, password, prefecture FROM user WHERE username = ?;", username)
	var user User
	if err := row.Scan(&user.id, &user.Username, &user.Password, &user.Prefecture); err != nil {
		log.Fatal(err)
	}
	return user
}
