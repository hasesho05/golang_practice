package db_client

import (
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/hasesho05/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitializeDBConnection() {
	dsn := "root:password@tcp(127.0.0.1:3306)/training?parseTime=true"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect DB:", err.Error())
	}

	db.AutoMigrate(&models.User{})

	DB = db
}
