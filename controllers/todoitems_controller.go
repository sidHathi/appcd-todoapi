package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"todo-api/models"
	"todo-api/services"

	"github.com/gin-gonic/gin"
)

// GetItem godoc
// @Summary Get a to-do item
// @Description Retrieve a specific to-do item by its ID
// @Tags todo-items
// @Produce  json
// @Param id path string true "Item ID"
// @Success 200 {object} models.TodoItem "Retrieved to-do item"
// @Failure 404 {string} string "Todo item not found"
// @Failure 400 {string} string "Error retrieving item"
// @Router /todo-item/{id} [get]
func GetItem(c *gin.Context) {
	itemId := c.Param("id")
	if itemId == "" {
		c.AbortWithStatus(404)
		return
	}

	fmt.Println("Item controller retrieving item with id " + itemId)
	item, err := services.GetItemById(itemId)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			c.AbortWithStatusJSON(404, "Todo item not found")
			return
		}
		c.AbortWithStatusJSON(400, err)
		return
	}
	c.JSON(http.StatusOK, *item)
}

// AddSubItem godoc
// @Summary Add a sub-item to a to-do item
// @Description Add a new sub-item to an existing to-do item
// @Tags todo-items
// @Accept  json
// @Produce  json
// @Param id path string true "Item ID"
// @Param item body models.CreateTodoItemData true "To-do sub-item data"
// @Success 201 {object} models.TodoItem "Sub-item added successfully"
// @Failure 404 {string} string "Todo item not found"
// @Failure 400 {string} string "Invalid input or error adding sub-item"
// @Router /todo-item/{id}/sub-item [post]
func AddSubItem(c *gin.Context) {
	itemId := c.Param("id")
	if itemId == "" {
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

	item, err := services.AddSubItem(itemId, body)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			c.AbortWithStatusJSON(404, "Todo item not found")
			return
		}
		c.AbortWithStatusJSON(400, err)
		return
	}
	c.JSON(http.StatusCreated, *item)
}

// UpdateItem godoc
// @Summary Update a to-do item
// @Description Update the details of an existing to-do item
// @Tags todo-items
// @Accept  json
// @Produce  json
// @Param id path string true "Item ID"
// @Param item body models.CreateTodoItemData true "Updated to-do item data"
// @Success 200 {string} string "Item updated successfully"
// @Failure 404 {string} string "Todo item not found"
// @Failure 400 {string} string "Invalid input or error updating item"
// @Router /todo-item/{id} [put]
func UpdateItem(c *gin.Context) {
	itemId := c.Param("id")
	if itemId == "" {
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

	err = services.UpdateItem(itemId, body)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			c.AbortWithStatusJSON(404, "Todo item not found")
			return
		}
		c.AbortWithStatusJSON(400, err)
		return
	}
	c.JSON(http.StatusOK, "Item updated successfully")
}

// SetItemCompletion godoc
// @Summary Set to-do item completion status
// @Description Update the completion status of a to-do item
// @Tags todo-items
// @Accept  json
// @Produce  json
// @Param id path string true "Item ID"
// @Param completion body models.TodoItemCompletionUpdate true "To-do item completion status update"
// @Success 200 {object} models.TodoItem "Completion status updated"
// @Failure 404 {string} string "Todo item not found"
// @Failure 400 {string} string "Invalid input or error updating completion status"
// @Router /todo-item/{id}/completion [put]
func SetItemCompletion(c *gin.Context) {
	itemId := c.Param("id")
	if itemId == "" {
		c.AbortWithStatus(404)
		return
	}

	body := models.TodoItemCompletionUpdate{}
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

	item, err := services.SetItemCompletion(itemId, body)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			c.AbortWithStatusJSON(404, "Todo item not found")
			return
		}
		c.AbortWithStatusJSON(400, err)
		return
	}
	c.JSON(http.StatusOK, *item)
}

// DeleteItem godoc
// @Summary Delete a to-do item
// @Description Delete an existing to-do item by its ID
// @Tags todo-items
// @Produce  json
// @Param id path string true "Item ID"
// @Success 200 {string} string "Item deleted successfully"
// @Failure 404 {string} string "Todo item not found"
// @Failure 400 {string} string "Error deleting item"
// @Router /todo-item/{id} [delete]
func DeleteItem(c *gin.Context) {
	itemId := c.Param("id")
	if itemId == "" {
		c.AbortWithStatus(404)
		return
	}

	err := services.DeleteItem(itemId)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			c.AbortWithStatusJSON(404, "Todo item not found")
			return
		}
		c.AbortWithStatusJSON(400, err)
		return
	}
	c.JSON(http.StatusOK, "Item deleted successfully")
}
