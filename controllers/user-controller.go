package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"todo-api/models"
	"todo-api/services"

	"github.com/gin-gonic/gin"
)

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
	} else {
		c.JSON(http.StatusCreated, *user)
	}
}

func GetUsers(c *gin.Context) {
	users, err := services.GetUsers()
	if err != nil {
		c.AbortWithStatusJSON(400, err.Error())
	}

	c.JSON(http.StatusOK, users)
}

func GetUserById(c *gin.Context) {
	userId := c.Param("id")
	user, err := services.GetUserById(userId)
	if err != nil {
		c.AbortWithStatusJSON(400, err.Error())
	}

	c.JSON(http.StatusOK, user)
}

func GetUserTodoLists(c *gin.Context) {
	userId := c.Param("id")
	listIds, err := services.GetUserListIds(userId)
	if err != nil {
		c.AbortWithStatusJSON(400, err.Error())
	}

	c.JSON(http.StatusOK, listIds)
}

func UpdateUser(c *gin.Context) {
	userId := c.Param("id")
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

func DeleteUser(c *gin.Context) {
	userId := c.Param("id")
	listIds, err := services.GetUserListIds(userId)
	if err != nil {
		c.AbortWithStatusJSON(400, err.Error())
	}

	c.JSON(http.StatusOK, listIds)
}
