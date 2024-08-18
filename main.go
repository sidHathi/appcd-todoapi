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
	router.GET("/users/:id", controllers.GetUserById)
	router.GET("/users/:id/todoListIds", controllers.GetUserTodoLists)
	router.PUT("/users/:id", controllers.UpdateUser)
	router.DELETE("/users/:id", controllers.DeleteUser)

	router.POST("/:userId/lists/", controllers.CreateList)
	router.GET("/:userId/lists/:listId", controllers.GetList)
	router.PUT("/:userId/lists/:listId", controllers.UpdateList)
	router.PUT("/:userId/lists/:listId/share", controllers.ShareList)
	router.POST("/:userId/lists/:listId/addItem", controllers.AddListItem)
	router.DELETE("/:userId/lists/:listId", controllers.DeleteList)

	router.GET("/items/:id", controllers.GetItem)
	router.POST("/items/:id/addSubItem", controllers.AddSubItem)
	router.PUT("/items/:id", controllers.UpdateItem)
	router.DELETE("/items/:id", controllers.DeleteItem)
	router.POST("/items/:id/addAttachment", controllers.AddAttachment)
	router.PUT("/items/attachments/:attId", controllers.UpdateAttachment)
	router.DELETE("/items/attachments/:attId", controllers.DeleteAttachment)

	router.Run("0.0.0.0:8000")
}
