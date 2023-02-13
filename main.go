package main

import (
	"go-sqlx-gin/controller"
	"go-sqlx-gin/db_client"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	db_client.InitializeDBConnection()

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:3000",
		},
		AllowHeaders: []string{
			"Origin",
			"Content-Length",
			"Content-Type",
		},
		AllowMethods: []string{
			"POST",
			"GET",
			"PUT",
			"DELETE",
			"OPTIONS",
		},
	}))

	r.GET("/user", controller.GetUserInfo)
	r.POST("/user/signup", controller.CreateUser)
	r.PUT("/user/signin", controller.Signin)
	r.PUT("/user/authorization", controller.Authorization)
	r.PUT("/user/changepassword", controller.ChangePassword)
	r.POST("/profile", controller.CreateProfile)
	r.POST("/profile/edit", controller.EditProfile)
	r.POST("/history", controller.CreateHistory)
	r.GET("/history/list", controller.GetHistory)

	if err := r.Run(":8000"); err != nil {
		log.Fatal(err)
	}
}
