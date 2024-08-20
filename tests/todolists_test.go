package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"todo-api/controllers"
	"todo-api/db"
	"todo-api/models"
	"todo-api/services"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCreateList(t *testing.T) {
	db.CleanTestDB()
	db.LoadTestData()
	t.Run("successful list creation", func(t *testing.T) {
		reqBody := `{"name": "New Test List"}`
		req, _ := http.NewRequest(http.MethodPost, "/users/tu_1/lists", bytes.NewBuffer([]byte(reqBody)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "userId", Value: "tu_1"}}
		c.Request = req

		controllers.CreateList(c)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response models.TodoList
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "New Test List", response.Name)

		// Verify that the list was actually created in the database
		list, err := services.GetListById("tu_1", response.Id)
		assert.NoError(t, err)
		assert.Equal(t, "New Test List", list.Name)
	})
}

func TestShareList(t *testing.T) {
	t.Run("successful list sharing", func(t *testing.T) {
		reqBody := `{"recipientId": "tu_2"}`
		req, _ := http.NewRequest(http.MethodPost, "/users/tu_1/lists/tl_1/share", bytes.NewBuffer([]byte(reqBody)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "userId", Value: "tu_1"}, {Key: "listId", Value: "tl_1"}}
		c.Request = req

		controllers.ShareList(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, `"List shared successfully"`, w.Body.String())

		// Verify that the list was actually shared in the database
		sharedList, err := services.GetListById("tu_2", "tl_1")
		assert.NoError(t, err)
		assert.Equal(t, "tl_1", sharedList.Id)
	})
}

func TestGetList(t *testing.T) {
	t.Run("successful list retrieval", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/users/tu_1/lists/tl_2", nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "userId", Value: "tu_1"}, {Key: "listId", Value: "tl_2"}}
		c.Request = req

		controllers.GetList(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response models.TodoList
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "tl_2", response.Id)
	})
}

func TestUpdateList(t *testing.T) {
	t.Run("successful list update", func(t *testing.T) {
		reqBody := `{"name": "Updated Test List"}`
		req, _ := http.NewRequest(http.MethodPut, "/users/tu_1/lists/tl_2", bytes.NewBuffer([]byte(reqBody)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "userId", Value: "tu_1"}, {Key: "listId", Value: "tl_2"}}
		c.Request = req

		controllers.UpdateList(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, `"List updated"`, w.Body.String())

		// Verify that the list was actually updated in the database
		list, err := services.GetListById("tu_1", "tl_2")
		assert.NoError(t, err)
		assert.Equal(t, "Updated Test List", list.Name)
	})
}

func TestAddListItem(t *testing.T) {
	t.Run("successful item addition to list", func(t *testing.T) {
		reqBody := `{"description": "New Test Item"}`
		req, _ := http.NewRequest(http.MethodPost, "/users/tu_1/lists/tl_1/items", bytes.NewBuffer([]byte(reqBody)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "userId", Value: "tu_1"}, {Key: "listId", Value: "tl_1"}}
		c.Request = req

		controllers.AddListItem(c)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response models.TodoItem
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "New Test Item", response.Description)

		// Verify that the item was actually added to the list in the database
		list, err := services.GetListById("tu_1", "tl_1")
		assert.NoError(t, err)
		assert.Contains(t, list.Items, response)
	})
}

func TestDeleteList(t *testing.T) {
	t.Run("successful list deletion", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodDelete, "/users/tu_2/lists/tl_3", nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "userId", Value: "tu_2"}, {Key: "listId", Value: "tl_3"}}
		c.Request = req

		controllers.DeleteList(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, `"List deleted"`, w.Body.String())

		// Verify that the list was actually deleted from the database
		_, err := services.GetListById("tu_2", "tl_3")
		assert.Error(t, err) // Expecting an error since the list should be deleted
	})
}
