package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

var db *sql.DB

type User struct {
	User_id   string
	User_name string
}

func connect() {
	// Construct the connection string
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))

	// Connect to the database
	db_local, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println("Connection error")
		panic(err)
	}
	db = db_local

	err = db.Ping()
	if err != nil {
		fmt.Println(connStr)
		fmt.Println(err.Error())
		panic(err)
	}
	fmt.Println("Database connection successful")
}

func helloWorld(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, "Hello world - successfully connected to postgres!")
}

func newUser(c *gin.Context) {
	body := User{}
	data, err := c.GetRawData()
	if err != nil {
		c.AbortWithStatusJSON(400, "No body provided for request")
		return
	}
	err = json.Unmarshal(data, &body)
	if err != nil {
		c.AbortWithStatusJSON(400, "Invalid input")
		return
	}

	_, err = db.Exec("insert into users(user_id,user_name) values ($1,$2);", body.User_id, body.User_name)
	if err != nil {
		fmt.Println(err)
		errString := fmt.Sprintf("Couldn't create the new user - err: %s", err.Error())
		c.AbortWithStatusJSON(400, errString)
	} else {
		c.JSON(http.StatusOK, "User is successfully created.")
	}
}

func getUsers(c *gin.Context) {
	data, err := db.Exec("select * from users;")
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatusJSON(400, "Get request failed")
	} else {
		c.JSON(http.StatusOK, data)
	}
}

func main() {
	connect()
	defer db.Close()

	router := gin.Default()
	router.GET("/", helloWorld)
	router.POST("/users/", newUser)
	router.GET("/users/", getUsers)

	router.Run("0.0.0.0:8000")
}
