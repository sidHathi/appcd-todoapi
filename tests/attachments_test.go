package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"todo-api/controllers"
	"todo-api/db"
	"todo-api/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAddAttachment(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db.CleanTestDB()
	db.LoadTestData()

	t.Run("successful attachment addition", func(t *testing.T) {
		body := `{"s3_url": "test_url", "file_type": "image"}`
		req, _ := http.NewRequest(http.MethodPost, "/items/ti_1/attachments", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "ti_1"}}
		c.Request = req

		controllers.AddAttachment(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response models.Attachment
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "test_url", response.S3Url)
	})

	t.Run("invalid input", func(t *testing.T) {
		body := `{"name": ""}`
		req, _ := http.NewRequest(http.MethodPost, "/items/1/attachments", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "1"}}
		c.Request = req

		controllers.AddAttachment(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestUpdateAttachment(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("successful attachment update", func(t *testing.T) {
		body := `{"s3_url": "Updated url", "file_type": "updated file type"}`
		req, _ := http.NewRequest(http.MethodPut, "/attachments/at_1", bytes.NewBuffer([]byte(body)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "at_1"}}
		c.Request = req

		controllers.UpdateAttachment(c)

		t.Log(w.Body)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("failed attachment update", func(t *testing.T) {
		body := `{"s3_url": "Failed Update"}`
		req, _ := http.NewRequest(http.MethodPut, "/attachments/4", bytes.NewBuffer([]byte(body)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "attId", Value: "4"}}
		c.Request = req

		controllers.UpdateAttachment(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestDeleteAttachment(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("successful deletion", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodDelete, "/attachments/at_2", nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "at_2"}}
		c.Request = req

		controllers.DeleteAttachment(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, `"Attachment deleted successfully"`, w.Body.String())
	})

	t.Run("failed deletion", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodDelete, "/attachments/4", nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "4"}}
		c.Request = req

		controllers.DeleteAttachment(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}
