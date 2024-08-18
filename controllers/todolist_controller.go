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
	}

	c.JSON(http.StatusCreated, *list)
}

func ShareList(c *gin.Context) {
	userId := c.Param("userId")
	listId := c.Param("listId")
	if userId == "" || listId == "" {
		c.AbortWithStatus(404)
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
	}
	c.JSON(http.StatusOK, "List shared successfully")
}

func GetList(c *gin.Context) {
	userId := c.Param("userId")
	listId := c.Param("listId")
	if userId == "" || listId == "" {
		c.AbortWithStatus(404)
	}

	list, err := services.GetListById(userId, listId)
	if err != nil {
		c.AbortWithStatusJSON(400, err)
	}
	c.JSON(http.StatusOK, list)
}

func UpdateList(c *gin.Context) {
	userId := c.Param("userId")
	listId := c.Param("listId")
	if userId == "" || listId == "" {
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

	err = services.UpdateListDetails(userId, listId, body)
	if err != nil {
		c.AbortWithStatusJSON(400, err)
	}
	c.JSON(http.StatusOK, "List updated")
}

func AddListItem(c *gin.Context) {
	userId := c.Param("userId")
	listId := c.Param("listId")
	if userId == "" || listId == "" {
		c.AbortWithStatus(404)
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
	}
	c.JSON(http.StatusCreated, item)
}

func DeleteList(c *gin.Context) {
	userId := c.Param("userId")
	listId := c.Param("listId")
	if userId == "" || listId == "" {
		c.AbortWithStatus(404)
	}

	err := services.DeleteList(userId, listId)
	if err != nil {
		c.AbortWithStatusJSON(400, err)
	}
	c.JSON(http.StatusOK, "List deleted")
}
