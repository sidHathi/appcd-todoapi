package controllers

import (
	"encoding/json"
	"net/http"
	"todo-api/models"
	"todo-api/services"

	"github.com/gin-gonic/gin"
)

func CreateList(c *gin.Context) {
	userId := c.Param("userId")
	if userId == "" {
		c.AbortWithStatus(404)
	}

	body := models.CreateTodoListData {}
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
	}

	c.JSON(http.StatusCreated, *list)
}