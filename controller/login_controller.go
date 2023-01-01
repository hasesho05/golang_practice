package controller

import (
	"go-sqlx-gin/db_client"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"
)

type User struct {
	Id         *int64  `json:"id"`
	Username   string  `json:"username"`
	Password   string  `json:"password"`
	Prefecture *string `json:"prefecture"`
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
			"content": err.Error(),
		})
		return
	}
	tmpPass, err := GetUser(resBody.Username)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status":  "ng",
			"message": "ユーザーが見つかりませんでした。",
		})
		return
	}
	dbPassword := tmpPass.Password
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

func ChangePassword(c *gin.Context) {
	var resBody struct {
		Username    string `json:"username"`
		Password    string `json:"password"`
		Newpassword string `json:"newpassword"`
	}
	if err := c.ShouldBindJSON(&resBody); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "Invalid request body",
			"content": err.Error(),
		})
		return
	}

	tmpPass, err := GetUser(resBody.Username)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status":  "ng",
			"message": "ユーザーが見つかりませんでした。",
		})
		return
	}
	dbPassword := tmpPass.Password
	formPassword := resBody.Password

	if err := CompareHashAndPassword(dbPassword, formPassword); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "ng",
			"message": "パスワードが違います。",
			"error":   err.Error(),
		})
		return
	}

	encryptedNewPassword, err := PasswordEncrypt(resBody.Newpassword)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "Invalid request body",
			"content": err.Error(),
		})
		return
	}

	res, err := db_client.DBClient.Exec("UPDATE user SET password = ? WHERE username = ?;",
		encryptedNewPassword,
		resBody.Username,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "ng",
			"content": err.Error(),
		})
		return
	}

	row, _ := res.RowsAffected()

	c.JSON(http.StatusAccepted, gin.H{
		"status":  "ok",
		"message": "パスワードが変更されました。",
		"row":     row,
	})
}

func GetUserInfo(c *gin.Context) {
	var resBody struct {
		Username string `json:"username"`
	}
	if err := c.ShouldBindJSON(&resBody); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "Invalid request body",
			"content": err.Error(),
		})
		return
	}
	user, err := GetUser(resBody.Username)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status":  "ng",
			"message": "ユーザーが見つかりません。",
		})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"status":     "ok",
		"username":   user.Username,
		"prefecture": user.Prefecture,
	})
}

func ChangePrefecture(c *gin.Context) {
	var resBody struct {
		Username   string `json:"username"`
		Prefecture string `json:"prefecture"`
	}
	if err := c.ShouldBindJSON(&resBody); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "Invalid request body",
			"content": err.Error(),
		})
		return
	}

	res, err := db_client.DBClient.Exec("UPDATE user SET prefecture = ? WHERE username = ?;",
		resBody.Prefecture,
		resBody.Username,
	)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status":  "ng",
			"message": "ユーザーが見つかりませんでした。",
			"content": err.Error(),
		})
		return
	}

	c.JSON(http.StatusUnprocessableEntity, gin.H{
		"status":   "ok",
		"message":  "都道府県情報を変更しました。",
		"response": res,
	})
}

func GetUser(username string) (User, error) {
	row := db_client.DBClient.QueryRow("SELECT id, username, password, prefecture FROM user WHERE username = ?;", username)
	var user User
	err := row.Scan(&user.Id, &user.Username, &user.Password, &user.Prefecture)
	return user, err
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
