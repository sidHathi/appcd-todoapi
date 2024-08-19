package controllers

import (
	"encoding/json"
	"net/http"
	"todo-api/models"
	"todo-api/services"

	"github.com/gin-gonic/gin"
)

// AddAttachment godoc
// @Summary Add an attachment to a to-do item
// @Description Add a new attachment to an existing to-do item
// @Tags attachments
// @Accept  json
// @Produce  json
// @Param id path string true "Item ID"
// @Param attachment body models.CreateAttachmentData true "Attachment data"
// @Success 200 {object} models.Attachment "Attachment added successfully"
// @Failure 404 {string} string "Todo item not found"
// @Failure 400 {string} string "Invalid input or error adding attachment"
// @Router /todo-items/:id/attachments [post]
func AddAttachment(c *gin.Context) {
	itemId := c.Param("id")
	if itemId == "" {
		c.AbortWithStatus(404)
		return
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
		return
	}
	c.JSON(http.StatusOK, att)
}

// UpdateAttachment godoc
// @Summary Update an existing attachment
// @Description Update the details of an existing attachment
// @Tags attachments
// @Accept  json
// @Produce  json
// @Param id path string true "Attachment ID"
// @Param attachment body models.CreateAttachmentData true "Updated attachment data"
// @Success 200 {string} string "Attachment updated successfully"
// @Failure 404 {string} string "Attachment not found"
// @Failure 400 {string} string "Invalid input or error updating attachment"
// @Router /attachments/{id} [put]
func UpdateAttachment(c *gin.Context) {
	attId := c.Param("id")
	if attId == "" {
		c.AbortWithStatus(404)
		return
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
		if err.Error() == "sql: no rows in result set" {
			c.AbortWithStatusJSON(404, "Attachment not found")
			return
		}
		c.AbortWithStatusJSON(400, err.Error())
		return
	}
	c.JSON(http.StatusOK, "Attachment updated successfully")
}

// DeleteAttachment godoc
// @Summary Delete an attachment
// @Description Delete an existing attachment by its ID
// @Tags attachments
// @Produce  json
// @Param id path string true "Attachment ID"
// @Success 200 {string} string "Attachment deleted successfully"
// @Failure 404 {string} string "Attachment not found"
// @Failure 400 {string} string "Error deleting attachment"
// @Router /attachments/{id} [delete]
func DeleteAttachment(c *gin.Context) {
	attId := c.Param("id")
	if attId == "" {
		c.AbortWithStatus(404)
		return
	}

	err := services.DeleteAttachment(attId)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			c.AbortWithStatusJSON(404, "Attachment not found")
			return
		}
		c.AbortWithStatusJSON(400, err)
		return
	}
	c.JSON(http.StatusOK, "Attachment deleted successfully")
}
