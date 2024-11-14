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

// ConnectDatabase establishes a connection to the PostgreSQL database using
// environment variables from a .env file. It constructs the connection string,
// opens the connection, assigns the database handle to the global variable DB,
// and prints a success message. If any error occurs, it prints an error message
// and panics. It also pings the database to ensure the connection is valid.
//
// Environment Variables:
// - DB_NAME: Database name
// - DB_PORT: Database port
// - DB_USER: Database user
// - DB_HOST: Database host
// - DB_PASSWORD: Database password
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
