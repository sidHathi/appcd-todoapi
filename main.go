package main

import (
	"fmt"
	"net/http"
	"todo-api/controllers"
	"todo-api/db"
	_ "todo-api/docs"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func helloWorld(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, "Hello world - successfully connected to postgres!")
}

func main() {
	// Load .env file for local runs
	err := godotenv.Load()
	if err != nil {
		fmt.Println("No .env file found, using default environment variables")
	}

	db.Connect()
	defer db.Db.Close()

	router := gin.Default()
	router.GET("/", helloWorld)

	router.POST("/users", controllers.CreateUser)
	router.GET("/users", controllers.GetUsers)
	router.GET("/users/:userId", controllers.GetUserById)
	router.GET("/users/:userId/todo-lists", controllers.GetUserTodoLists)
	router.PUT("/users/:userId", controllers.UpdateUser)
	router.DELETE("/users/:userId", controllers.DeleteUser)

	router.POST("/users/:userId/todo-lists", controllers.CreateList)
	router.GET("/users/:userId/todo-lists/:listId", controllers.GetList)
	router.PUT("/users/:userId/todo-lists/:listId", controllers.UpdateList)
	router.PUT("/users/:userId/todo-lists/:listId/share", controllers.ShareList)
	router.POST("/users/:userId/todo-lists/:listId/items", controllers.AddListItem)
	router.DELETE("/users/:userId/todo-lists/:listId", controllers.DeleteList)

	router.GET("/todo-items/:id", controllers.GetItem)
	router.POST("/todo-items/:id", controllers.AddSubItem)
	router.PUT("/todo-items/:id", controllers.UpdateItem)
	router.PUT("/todo-items/:id/completion", controllers.SetItemCompletion)
	router.DELETE("/todo-items/:id", controllers.DeleteItem)
	router.POST("/todo-items/:id/attachments", controllers.AddAttachment)

	router.PUT("/attachments/:id", controllers.UpdateAttachment)
	router.DELETE("/attachments/:id", controllers.DeleteAttachment)

	// Serve Swagger docs
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.Run("0.0.0.0:8000")
}
