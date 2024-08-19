package controllers

import (
	"encoding/json"
	"net/http"
	"todo-api/models"
	"todo-api/services"

	"github.com/gin-gonic/gin"
)

// CreateList godoc
// @Summary Create a new to-do list
// @Description Create a new to-do list for a specific user
// @Tags todo-lists
// @Accept  json
// @Produce  json
// @Param userId path string true "User ID"
// @Param list body models.CreateTodoListData true "To-do list data"
// @Success 201 {object} models.TodoList "List created successfully"
// @Failure 400 {string} string "Invalid input or error creating list"
// @Router /users/{userId}/todo-lists [post]
func CreateList(c *gin.Context) {
	userId := c.Param("userId")
	if userId == "" {
		c.AbortWithStatus(404)
	}

	body := models.CreateTodoListData{}
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

	list, err := services.CreateNewList(body, userId)
	if err != nil {
		c.AbortWithStatusJSON(400, err)
		return
	}

	c.JSON(http.StatusCreated, *list)
}

// ShareList godoc
// @Summary Share a to-do list with another user
// @Description Share a specific to-do list with another user by providing their user ID
// @Tags todo-lists
// @Accept  json
// @Produce  json
// @Param userId path string true "User ID"
// @Param listId path string true "List ID"
// @Param share body models.ShareListData true "Share list data"
// @Success 200 {string} string "List shared successfully"
// @Failure 400 {string} string "Invalid input or error sharing list"
// @Router /users/{userId}/todo-lists/{listId}/share [post]
func ShareList(c *gin.Context) {
	userId := c.Param("userId")
	listId := c.Param("listId")
	if userId == "" || listId == "" {
		c.AbortWithStatus(404)
		return
	}

	body := models.ShareListData{}
	data, err := c.GetRawData()
	if err != nil {
		c.AbortWithStatusJSON(400, "No body provided for request")
		return
	}
	err = json.Unmarshal(data, &body)
	if err != nil || body.RecipientId == "" {
		c.AbortWithStatusJSON(400, "Invalid input")
		return
	}

	err = services.ShareList(userId, listId, body.RecipientId)
	if err != nil {
		c.AbortWithStatusJSON(400, err)
		return
	}
	c.JSON(http.StatusOK, "List shared successfully")
}

// GetList godoc
// @Summary Get a specific to-do list
// @Description Retrieve a specific to-do list by its ID and the user's ID
// @Tags todo-lists
// @Produce  json
// @Param userId path string true "User ID"
// @Param listId path string true "List ID"
// @Success 200 {object} models.TodoList "Retrieved to-do list"
// @Failure 400 {string} string "Error retrieving list"
// @Router /users/{userId}/todo-lists/{listId} [get]
func GetList(c *gin.Context) {
	userId := c.Param("userId")
	listId := c.Param("listId")
	if userId == "" || listId == "" {
		c.AbortWithStatus(404)
		return
	}

	list, err := services.GetListById(userId, listId)
	if err != nil {
		c.AbortWithStatusJSON(400, err)
		return
	}
	c.JSON(http.StatusOK, *list)
}

// UpdateList godoc
// @Summary Update a to-do list
// @Description Update the details of an existing to-do list
// @Tags todo-lists
// @Accept  json
// @Produce  json
// @Param userId path string true "User ID"
// @Param listId path string true "List ID"
// @Param list body models.CreateTodoListData true "Updated to-do list data"
// @Success 200 {string} string "List updated"
// @Failure 400 {string} string "Invalid input or error updating list"
// @Router /users/{userId}/todo-lists/{listId} [put]
func UpdateList(c *gin.Context) {
	userId := c.Param("userId")
	listId := c.Param("listId")
	if userId == "" || listId == "" {
		c.AbortWithStatus(404)
		return
	}

	body := models.CreateTodoListData{}
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

	err = services.UpdateListDetails(userId, listId, body)
	if err != nil {
		c.AbortWithStatusJSON(400, err)
		return
	}
	c.JSON(http.StatusOK, "List updated")
}

// AddListItem godoc
// @Summary Add an item to a to-do list
// @Description Add a new item to an existing to-do list
// @Tags todo-lists
// @Accept  json
// @Produce  json
// @Param userId path string true "User ID"
// @Param listId path string true "List ID"
// @Param item body models.CreateTodoItemData true "To-do item data"
// @Success 201 {object} models.TodoItem "Item added successfully"
// @Failure 400 {string} string "Invalid input or error adding item"
// @Router /users/{userId}/todo-lists/{listId}/items [post]
func AddListItem(c *gin.Context) {
	userId := c.Param("userId")
	listId := c.Param("listId")
	if userId == "" || listId == "" {
		c.AbortWithStatus(404)
		return
	}

	body := models.CreateTodoItemData{}
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

	item, err := services.AddListItem(userId, listId, body)
	if err != nil {
		c.AbortWithStatusJSON(400, err)
		return
	}
	c.JSON(http.StatusCreated, *item)
}

// DeleteList godoc
// @Summary Delete a to-do list
// @Description Delete an existing to-do list by its ID
// @Tags todo-lists
// @Produce  json
// @Param userId path string true "User ID"
// @Param listId path string true "List ID"
// @Success 200 {string} string "List deleted"
// @Failure 400 {string} string "Error deleting list"
// @Router /users/{userId}/todo-lists/{listId} [delete]
func DeleteList(c *gin.Context) {
	userId := c.Param("userId")
	listId := c.Param("listId")
	if userId == "" || listId == "" {
		c.AbortWithStatus(404)
		return
	}

	err := services.DeleteList(userId, listId)
	if err != nil {
		c.AbortWithStatusJSON(400, err.Error())
		return
	}
	c.JSON(http.StatusOK, "List deleted")
}
