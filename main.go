package main

import (
	"go-sqlx-gin/controller"
	"go-sqlx-gin/db_client"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	db_client.InitializeDBConnection()

	r := gin.Default()

	r.POST("/", controller.CreatePost)
	r.GET("/", controller.GetPosts)
	r.GET("/:id", controller.GetPost)
	r.POST("/user/signup", controller.CreateUser)
	r.GET("/user/login", controller.Login)

	if err := r.Run(":5000"); err != nil {
		log.Fatal(err)
	}
}
