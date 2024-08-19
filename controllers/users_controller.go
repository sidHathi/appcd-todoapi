package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"todo-api/models"
	"todo-api/services"

	"github.com/gin-gonic/gin"
)

// CreateUser godoc
// @Summary Create a new user
// @Description Create a new user with the provided information
// @Tags users
// @Accept  json
// @Produce  json
// @Param user body models.CreateUserModel true "User data"
// @Success 201 {object} models.User
// @Failure 400 {string} string "Invalid input or error creating user"
// @Router /users [post]
func CreateUser(c *gin.Context) {
	body := models.CreateUserModel{}
	data, err := c.GetRawData()
	if err != nil {
		c.AbortWithStatusJSON(400, "No body provided for request")
		return
	}
	err = json.Unmarshal(data, &body)
	if err != nil || body.Name == "" {
		c.AbortWithStatusJSON(400, "Invalid input")
		return
	}

	user, err := services.CreateUser(body)
	if err != nil {
		errString := fmt.Sprintf("Couldn't create the new user - err: %s", err.Error())
		c.AbortWithStatusJSON(400, errString)
		return
	} else {
		c.JSON(http.StatusCreated, *user)
	}
}

// GetUsers godoc
// @Summary Get a list of users
// @Description Retrieve a list of all users
// @Tags users
// @Produce  json
// @Success 200 {array} models.User
// @Failure 400 {string} string "Error retrieving users"
// @Router /users [get]
func GetUsers(c *gin.Context) {
	users, err := services.GetUsers()
	if err != nil {
		c.AbortWithStatusJSON(400, err.Error())
		return
	}

	c.JSON(http.StatusOK, users)
}

// GetUserById godoc
// @Summary Get a user by ID
// @Description Retrieve a user by their ID
// @Tags users
// @Produce  json
// @Param userId path string true "User ID"
// @Success 200 {object} models.User
// @Failure 404 {string} string "User not found"
// @Failure 400 {string} string "Error retrieving user"
// @Router /users/{userId} [get]
func GetUserById(c *gin.Context) {
	userId := c.Param("userId")
	user, err := services.GetUserById(userId)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			c.AbortWithStatusJSON(404, "User not found")
			return
		}
		c.AbortWithStatusJSON(400, err.Error())
		return
	}

	c.JSON(http.StatusOK, user)
}

// GetUserTodoLists godoc
// @Summary Get user's to-do lists
// @Description Retrieve a list of to-do lists associated with the specified user
// @Tags users
// @Produce  json
// @Param userId path string true "User ID"
// @Success 200 {array} string "List of to-do list IDs"
// @Failure 400 {string} string "Error retrieving to-do lists"
// @Router /users/{userId}/todo-lists [get]
func GetUserTodoLists(c *gin.Context) {
	userId := c.Param("userId")
	listIds, err := services.GetUserListIds(userId)
	if err != nil {
		c.AbortWithStatusJSON(400, err.Error())
		return
	}

	c.JSON(http.StatusOK, listIds)
}

// UpdateUser godoc
// @Summary Update a user
// @Description Update the information of an existing user
// @Tags users
// @Accept  json
// @Produce  json
// @Param userId path string true "User ID"
// @Param user body models.CreateUserModel true "Updated user data"
// @Success 200 {object} models.CreateUserModel
// @Failure 400 {string} string "Invalid input or error updating user"
// @Router /users/{userId} [put]
func UpdateUser(c *gin.Context) {
	userId := c.Param("userId")
	body := models.CreateUserModel{}
	data, err := c.GetRawData()
	if err != nil {
		c.AbortWithStatusJSON(400, "No body provided for request")
		return
	}

	err = json.Unmarshal(data, &body)
	if err != nil || body.Name == "" {
		c.AbortWithStatusJSON(400, "Invalid input")
		return
	}
	err = services.UpdateUser(userId, body)
	if err != nil {
		c.AbortWithStatusJSON(400, err.Error())
	}

	c.JSON(http.StatusOK, body)
}

// DeleteUser godoc
// @Summary Delete a user
// @Description Delete an existing user by their ID
// @Tags users
// @Produce  json
// @Param userId path string true "User ID"
// @Success 200 {string} string "User deleted successfully"
// @Failure 400 {string} string "Error deleting user"
// @Router /users/{userId} [delete]
func DeleteUser(c *gin.Context) {
	userId := c.Param("userId")
	err := services.DeleteUser(userId)
	if err != nil {
		c.AbortWithStatusJSON(400, err.Error())
		return
	}

	c.JSON(http.StatusOK, "User deleted successfully")
}
