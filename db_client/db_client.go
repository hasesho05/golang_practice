package db_client

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var DBClient *sql.DB

func InitializeDBConnection() {
	db, err := sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/training?parseTime=true")
	if err != nil {
		log.Fatal("failed to connect DB:", err.Error())
	}
	err = db.Ping()
	if err != nil {
		log.Fatal("failed to connect DB:", err.Error())
	}

	DBClient = db
}
