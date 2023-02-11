package controller

import (
	"fmt"
	"go-sqlx-gin/db_client"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"
)

type User struct {
	Id        *int64  `json:"id"`
	Username  *string `json:"username"`
	Email     *string `json:"email"`
	Password  string  `json:"password"`
	Token     string  `json:"token"`
	Icon      *string `json:"icon"`
	CreatedAt *string `json:"created_at"`
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
	fmt.Println(resBody)

	if isExistUser(resBody.Token) {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "User already exists",
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

	_, err = db_client.DBClient.Exec("INSERT INTO user (username, email, password, icon, token) VALUES (?, ?, ?, ?, ?);",
		resBody.Username,
		resBody.Email,
		encryptedPassword,
		resBody.Icon,
		resBody.Token,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "ng",
			"content": err.Error(),
			"message": "Invalid request body",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status": "ok",
		"error":  false,
	})
}

func Signin(c *gin.Context) {
	var resBody User
	if err := c.ShouldBindJSON(&resBody); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "Invalid request body",
			"content": err.Error(),
		})
	}
	user, err := GetUser(*resBody.Email)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status":  "ng",
			"message": "ユーザーが見つかりませんでした。",
		})
	}
	dbPassword := user.Password
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
		"data":    user,
	})
}

func Authorization(c *gin.Context) {
	var resBody struct {
		Token string `json:"token"`
	}
	if err := c.ShouldBindJSON(&resBody); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "Invalid request body",
			"content": err.Error(),
		})
		return
	}
	user, err := GetUserByToken(resBody.Token)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status":  "ng",
			"message": "ユーザーが見つかりませんでした。",
		})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"status":  "ok",
		"message": "ログインしました。",
		"data":    user,
	})
}

func GetUserByToken(token string) (User, error) {
	var user User
	err := db_client.DBClient.QueryRow("SELECT * FROM user WHERE token = ?;", token).Scan(
		&user.Id,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.Token,
		&user.Icon,
		&user.CreatedAt,
	)
	if err != nil {
		log.Println(err)
		return user, err
	}
	return user, nil
}

func ChangePassword(c *gin.Context) {
	var resBody struct {
		Email       *string `json:"email"`
		Password    string  `json:"password"`
		Newpassword string  `json:"newpassword"`
	}
	if err := c.ShouldBindJSON(&resBody); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "Invalid request body",
			"content": err.Error(),
		})
		return
	}

	user, err := GetUser(*resBody.Email)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status":  "ng",
			"message": "ユーザーが見つかりませんでした。",
		})
		return
	}
	dbPassword := user.Password
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

	_, err = db_client.DBClient.Exec("UPDATE user SET password = ? WHERE email = ?;",
		encryptedNewPassword,
		resBody.Email,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "ng",
			"content": err.Error(),
		})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"status":  "ok",
		"message": "パスワードが変更されました。",
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
		"status":   "ok",
		"username": user.Username,
		"icon":     user.Icon,
	})
}

func GetUser(email string) (User, error) {
	row := db_client.DBClient.QueryRow("SELECT id, username, email, password, icon, token FROM user WHERE email = ?;", email)
	var user User
	err := row.Scan(&user.Id, &user.Username, &user.Email, &user.Password, &user.Icon, &user.Token)
	if err != nil {
		log.Fatal(err.Error())
	}
	return user, err
}

func isExistUser(token string) bool {
	row := db_client.DBClient.QueryRow("SELECT id, username, password, icon, token FROM user WHERE token = ?;", token)
	return row == nil
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
