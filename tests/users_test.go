package tests

import (
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

func TestCreateUser(t *testing.T) {
	db.CleanTestDB()
	db.LoadTestData()
	gin.SetMode(gin.TestMode)
	
	t.Run("successful creation", func(t *testing.T) {
		body := `{"name": "John Doe"}`
		req, _ := http.NewRequest(http.MethodPost, "/users", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		controllers.CreateUser(c)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response models.User
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "John Doe", response.Name)
	})

	t.Run("failed creation due to invalid input", func(t *testing.T) {
		body := `{"name": ""}`
		req, _ := http.NewRequest(http.MethodPost, "/users", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		controllers.CreateUser(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestGetUserById(t *testing.T) {
	db.CleanTestDB()
	db.LoadTestData()
	gin.SetMode(gin.TestMode)
	
	t.Run("user found", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/users/tu_1", nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "userId", Value: "tu_1"}}
		c.Request = req

		controllers.GetUserById(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response models.User
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Test User 1", response.Name)
	})

	t.Run("user not found", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/users/invalid", nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "userId", Value: "invalid"}}
		c.Request = req

		controllers.GetUserById(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestUpdateUser(t *testing.T) {
	db.CleanTestDB()
	db.LoadTestData()
	gin.SetMode(gin.TestMode)
	
	t.Run("successful update", func(t *testing.T) {
		body := `{"name": "Jane Doe"}`
		req, _ := http.NewRequest(http.MethodPut, "/users/test_id", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "userId", Value: "test_id"}}
		c.Request = req

		controllers.UpdateUser(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response models.CreateUserModel
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Jane Doe", response.Name)
	})

	t.Run("failed update due to invalid input", func(t *testing.T) {
		body := `{"name": ""}`
		req, _ := http.NewRequest(http.MethodPut, "/users/test_id", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "userId", Value: "test_id"}}
		c.Request = req

		controllers.UpdateUser(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestDeleteUser(t *testing.T) {
	db.CleanTestDB()
	db.LoadTestData()
	gin.SetMode(gin.TestMode)
	
	t.Run("successful deletion", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodDelete, "/users/test_id", nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "userId", Value: "test_id"}}
		c.Request = req

		controllers.DeleteUser(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, `"User deleted successfully"`, w.Body.String())
	})

	t.Run("failed deletion", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodDelete, "/users/2", nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "userId", Value: "2"}}
		c.Request = req

		controllers.DeleteUser(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
