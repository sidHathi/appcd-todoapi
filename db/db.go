package db

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

func ConnectForTesting() {
	// Construct the connection string
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("TEST_DB_HOST"), os.Getenv("TEST_DB_PORT"), os.Getenv("TEST_DB_USER"), os.Getenv("TEST_DB_PASSWORD"), os.Getenv("TEST_DB_NAME"))

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

	CleanTestDB()
}

func CleanTestDB() {
	Db.Exec(`
		delete from users;
		delete from todo_lists;
		delete from attachments;
		delete from todo_items;
		delete from user_todo_lists;
	`)
}

func LoadTestData() {
	// init test db:
	Db.Exec("insert into users (id, name) values ($1, $2);", "test_id", "test_name")
	Db.Exec("insert into users (id, name) values ($1, $2);", "tu_1", "Test User 1")
	Db.Exec("insert into users (id, name) values ($1, $2);", "tu_2", "Test User 2")

	Db.Exec("insert into todo_items (id, list_id, description, complete, parent_id) values ($1, $2, $3, $4, $5);", "ti_1", "tl_1", "Test Item 1", false, "")
	Db.Exec("insert into todo_items (id, list_id, description, complete, parent_id) values ($1, $2, $3, $4, $5);", "ti_2", "tl_1", "Test Item 2", false, "")

	Db.Exec("insert into attachments (id, list_id, item_id, s3_url, file_type) values ($1, $2, $3, $4, $5);", "at_1", "tl_1", "ti_1", "Att_url", "image")
	Db.Exec("insert into attachments (id, list_id, item_id, s3_url, file_type) values ($1, $2, $3, $4, $5);", "at_2", "tl_1", "ti_1", "Att_url_2", "image")
	Db.Exec("insert into attachments (id, list_id, item_id, s3_url, file_type) values ($1, $2, $3, $4, $5);", "at_3", "tl_1", "ti_1", "Att_url", "image")

	Db.Exec("insert into todo_lists (id, name, created_by) values ($1, $2, $3);", "tl_1", "Test List 1", "tu_1")
	Db.Exec("insert into todo_lists (id, name, created_by) values ($1, $2, $3);", "tl_2", "Test List 2", "tu_1")
	Db.Exec("insert into todo_lists (id, name, created_by) values ($1, $2, $3);", "tl_3", "Test List 3", "tu_2")
	Db.Exec("insert into todo_lists (id, name, created_by) values ($1, $2, $3);", "tl_4", "Integration Testing List", "tu_1")

	Db.Exec("insert into user_todo_lists (list_id, user_id) values ($1, $2);", "tl_1", "tu_1")
	Db.Exec("insert into user_todo_lists (list_id, user_id) values ($1, $2);", "tl_2", "tu_1")
	Db.Exec("insert into user_todo_lists (list_id, user_id) values ($1, $2);", "tl_3", "tu_2")
	Db.Exec("insert into user_todo_lists (list_id, user_id) values ($1, $2);", "tl_4", "tu_1")
}
