package database

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var DB *sql.DB

func GetDB() *sql.DB {
	return DB
}

func ConnectDatabase() {

	err := godotenv.Load()

	if err != nil {
		fmt.Println("Error on .env file.", err)
	}

	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	host := os.Getenv("DB_HOST")
	pass := os.Getenv("DB_PASSWORD")

	psqlSetup := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		host, port, user, dbname, pass)

	db, errSql := sql.Open("postgres", psqlSetup)

	if errSql != nil {
		fmt.Println("Error connecting to the database ", err)
		panic(err)
	} else {
		DB = db
		fmt.Println("Successfully connected to database!")
	}

	// Test connection
	err = DB.Ping()
	if err != nil {
		fmt.Println("Error on .env file.", err)
	}

}
