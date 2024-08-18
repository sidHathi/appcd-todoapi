package main

import (
	"net/http"
	"todo-api/controllers"
	"todo-api/db"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func helloWorld(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, "Hello world - successfully connected to postgres!")
}

func main() {
	db.Connect()
	defer db.Db.Close()

	router := gin.Default()
	router.GET("/", helloWorld)
	router.POST("/users/", controllers.CreateUser)
	router.GET("/users/", controllers.GetUsers)

	router.Run("0.0.0.0:8000")
}
