package controllers

import (
	"encoding/json"
	"net/http"
	"todo-api/models"
	"todo-api/services"

	"github.com/gin-gonic/gin"
)

func GetItem(c *gin.Context) {
	itemId := c.Param("id")
	if itemId == "" {
		c.AbortWithStatus(404)
	}

	item, err := services.GetItemById(itemId)
	if err != nil {
		c.AbortWithStatusJSON(400, err)
	}
	c.JSON(http.StatusOK, item)
}

func AddSubItem(c *gin.Context) {
	itemId := c.Param("id")
	if itemId == "" {
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

	item, err := services.AddSubItem(itemId, body)
	if err != nil {
		c.AbortWithStatusJSON(400, err)
	}
	c.JSON(http.StatusOK, item)
}

func UpdateItem(c *gin.Context) {
	itemId := c.Param("id")
	if itemId == "" {
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

	err = services.UpdateItem(itemId, body)
	if err != nil {
		c.AbortWithStatusJSON(400, err)
	}
	c.JSON(http.StatusOK, "Item updated successfully")
}

func AddAttachment(c *gin.Context) {
	itemId := c.Param("id")
	if itemId == "" {
		c.AbortWithStatus(404)
	}

	body := models.CreateAttachmentData{}
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

	att, err := services.AddAttachment(itemId, body)
	if err != nil {
		c.AbortWithStatusJSON(400, err)
	}
	c.JSON(http.StatusOK, att)
}

func UpdateAttachment(c *gin.Context) {
	attId := c.Param("attId")
	if attId == "" {
		c.AbortWithStatus(404)
	}

	body := models.CreateAttachmentData{}
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

	err = services.UpdateAttachment(attId, body)
	if err != nil {
		c.AbortWithStatusJSON(400, err)
	}
	c.JSON(http.StatusOK, "Attachment updated successfully")
}

func DeleteAttachment(c *gin.Context) {
	attId := c.Param("attId")
	if attId == "" {
		c.AbortWithStatus(404)
	}

	err := services.DeleteAttachment(attId)
	if err != nil {
		c.AbortWithStatusJSON(400, err)
	}
	c.JSON(http.StatusOK, "Attachment deleted successfully")
}

func DeleteItem(c *gin.Context) {
	itemId := c.Param("id")
	if itemId == "" {
		c.AbortWithStatus(404)
	}

	err:= services.DeleteItem(itemId)
	if err != nil {
		c.AbortWithStatusJSON(400, err)
	}
	c.JSON(http.StatusOK, "Item deleted successfully")
}