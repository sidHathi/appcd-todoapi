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

func TestGetItem(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db.CleanTestDB()
	db.LoadTestData()

	t.Run("item found", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/items/ti_1", nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "ti_1"}}
		c.Request = req
		t.Log("Executing get request!")

		controllers.GetItem(c)

		assert.Equal(t, http.StatusOK, w.Code)
		t.Log("Item found executed get with status 200!")
		var response models.TodoItemNoNest
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Test Item 1", response.Description)
	})

	t.Run("item not found", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/items/2", nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "2"}}
		c.Request = req

		controllers.GetItem(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestAddSubItem(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("successful sub-item creation", func(t *testing.T) {
		body := `{"description": "Sub Item 1"}`
		req, _ := http.NewRequest(http.MethodPost, "/items/ti_1/addSubItem", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "ti_1"}}
		c.Request = req

		controllers.AddSubItem(c)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response models.TodoItem
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Sub Item 1", response.Description)
	})

	t.Run("invalid input", func(t *testing.T) {
		body := `{"name": ""}`
		req, _ := http.NewRequest(http.MethodPost, "/items/2/addSubItem", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "2"}}
		c.Request = req

		controllers.AddSubItem(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestUpdateItem(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("successful update", func(t *testing.T) {
		body := `{"description": "Updated Item Description"}`
		req, _ := http.NewRequest(http.MethodPut, "/items/ti_1", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "ti_1"}}
		c.Request = req

		controllers.UpdateItem(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("failed update", func(t *testing.T) {
		body := `{"description": "Failed Update"}`
		req, _ := http.NewRequest(http.MethodPut, "/items/invalid", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "invalid"}}
		c.Request = req

		controllers.UpdateItem(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestSetItemCompletion(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("successful completion update", func(t *testing.T) {
		body := `{"complete": true}`
		req, _ := http.NewRequest(http.MethodPut, "/items/ti_1/completion", bytes.NewBuffer([]byte(body)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "ti_1"}}
		c.Request = req

		controllers.SetItemCompletion(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response models.TodoItem
		err := json.Unmarshal(w.Body.Bytes(), &response)
		t.Log(response)
		assert.NoError(t, err)
		assert.True(t, response.Complete)
	})

	t.Run("invalid input", func(t *testing.T) {
		body := `{"complete": "invalid"}`
		req, _ := http.NewRequest(http.MethodPut, "/items/ti_1/completion", bytes.NewBuffer([]byte(body)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "1"}}
		c.Request = req

		controllers.SetItemCompletion(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestDeleteItem(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("successful deletion", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodDelete, "/items/ti_2", nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "ti_2"}}
		c.Request = req

		controllers.DeleteItem(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, `"Item deleted successfully"`, w.Body.String())
	})

	t.Run("failed deletion", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodDelete, "/items/2", nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "2"}}
		c.Request = req

		controllers.DeleteItem(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}
