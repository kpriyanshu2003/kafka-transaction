package config

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

type Dbinstance struct {
	Db *sql.DB
}

var DB Dbinstance

func SetupDatabase() {
	db_url := os.Getenv("DB_URL")
	if db_url == "" {
		log.Fatal("DB_URL is missing")
	}

	conn, err := sql.Open("postgres", db_url)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to Database")
	DB = Dbinstance{Db: conn}
}
