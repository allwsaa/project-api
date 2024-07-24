package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var datab *sql.DB

func SetupDB() {
	var err error

	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_PORT"), os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_DB"))
	datab, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal("Error opening database connection: ", err)
	}

	err = datab.Ping()
	if err != nil {
		log.Fatal("Error pinging database: ", err)
	}

	log.Println("Database connected successfully")
}

func GetDB() *sql.DB {
	return datab
}
