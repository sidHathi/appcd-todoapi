package db;

import (
	"database/sql"
	"fmt"
	"os"
)

var Db *sql.DB

func Connect() {
	// Construct the connection string
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))

	// Connect to the database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println("Connection error")
		panic(err)
	}
	Db = db

	err = db.Ping()
	if err != nil {
		fmt.Println(connStr)
		fmt.Println(err.Error())
		panic(err)
	}
	fmt.Println("Database connection successful")
}
