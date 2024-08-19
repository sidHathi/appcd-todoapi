package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
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

func TestMultiNestedItemLists(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db.CleanTestDB()
	db.LoadTestData()

	userId := "tu_1"
	listId := "tl_4"
	createItemPath := fmt.Sprintf("/users/%s/lists/%s/items", userId, listId)
	getListPath := fmt.Sprintf("/users/%s/lists/%s", userId, listId)
	var item1Id string
	var item2Id string
	var item3Id string

	t.Run("successful creation of multi tiered list item structure", func (t *testing.T) {
		// Assuming userId 1 and listId 2 exist in your test database
		reqBody := `{"description": "Test Item 1"}`
		req, _ := http.NewRequest(http.MethodPost, createItemPath, bytes.NewBuffer([]byte(reqBody)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "userId", Value: userId}, {Key: "listId", Value: listId}}
		c.Request = req

		controllers.AddListItem(c)
		var response models.TodoItem
		assert.Equal(t, http.StatusCreated, w.Code)
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		item1Id = response.Id

		subItemPath := fmt.Sprintf("/items/%s", item1Id)
		reqBody = `{"description": "Test Item 2"}`
		req, _ = http.NewRequest(http.MethodPost, subItemPath, bytes.NewBuffer([]byte(reqBody)))
		req.Header.Set("Content-Type", "application/json")

		w = httptest.NewRecorder()
		c, _ = gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: item1Id}}
		c.Request = req

		controllers.AddSubItem(c)
		assert.Equal(t, http.StatusCreated, w.Code)
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		item2Id = response.Id

		subItemPath = fmt.Sprintf("/items/%s", item2Id)
		reqBody = `{"description": "Test Item 3"}`
		req, _ = http.NewRequest(http.MethodPost, subItemPath, bytes.NewBuffer([]byte(reqBody)))
		req.Header.Set("Content-Type", "application/json")

		w = httptest.NewRecorder()
		c, _ = gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: item2Id}}
		c.Request = req

		controllers.AddSubItem(c)
		assert.Equal(t, http.StatusCreated, w.Code)
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		item3Id = response.Id

		req, _ = http.NewRequest(http.MethodGet, getListPath, nil)
		w = httptest.NewRecorder()
		c, _ = gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "userId", Value: userId}, {Key: "listId", Value: listId}}
		c.Request = req

		controllers.GetList(c)
		assert.Equal(t, http.StatusOK, w.Code)

		var list models.TodoList
		err = json.Unmarshal(w.Body.Bytes(), &list)
		assert.NoError(t, err)
		t.Log(list)

		assert.NotEmpty(t, list.Items)
		assert.Equal(t, list.Items[0].Id, item1Id)
		assert.NotEmpty(t, list.Items[0].SubItems)
		assert.Equal(t, list.Items[0].SubItems[0].Id, item2Id)
		assert.NotEmpty(t, list.Items[0].SubItems[0].SubItems)
		assert.Equal(t, list.Items[0].SubItems[0].SubItems[0].Id, item3Id)
	})

	t.Run("successful completion of layered list items", func (t *testing.T) {
		completionPath := fmt.Sprintf("/items/%s/completion", item1Id)
		body := `{"complete": true}`
		req, _ := http.NewRequest(http.MethodPut, completionPath, bytes.NewBuffer([]byte(body)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: item1Id}}
		c.Request = req

		controllers.SetItemCompletion(c)
		assert.Equal(t, http.StatusOK, w.Code)

		req, _ = http.NewRequest(http.MethodGet, getListPath, nil)
		w = httptest.NewRecorder()
		c, _ = gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "userId", Value: userId}, {Key: "listId", Value: listId}}
		c.Request = req

		controllers.GetList(c)
		assert.Equal(t, http.StatusOK, w.Code)

		var list models.TodoList
		err := json.Unmarshal(w.Body.Bytes(), &list)
		assert.NoError(t, err)
		assert.NotEmpty(t, list.Items)
		assert.True(t, list.Items[0].Complete)
		assert.NotEmpty(t, list.Items[0].SubItems)
		assert.True(t, list.Items[0].SubItems[0].Complete)
		assert.NotEmpty(t, list.Items[0].SubItems[0].SubItems)
		assert.True(t, list.Items[0].SubItems[0].SubItems[0].Complete)

		completionPath = fmt.Sprintf("/items/%s/completion", item2Id)
		body = `{"complete": false}`
		req, _ = http.NewRequest(http.MethodPut, completionPath, bytes.NewBuffer([]byte(body)))
		req.Header.Set("Content-Type", "application/json")

		w = httptest.NewRecorder()
		c, _ = gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: item2Id}}
		c.Request = req

		controllers.SetItemCompletion(c)
		assert.Equal(t, http.StatusOK, w.Code)

		req, _ = http.NewRequest(http.MethodGet, getListPath, nil)
		w = httptest.NewRecorder()
		c, _ = gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "userId", Value: userId}, {Key: "listId", Value: listId}}
		c.Request = req

		controllers.GetList(c)
		assert.Equal(t, http.StatusOK, w.Code)
		err = json.Unmarshal(w.Body.Bytes(), &list)
		assert.NoError(t, err)
		assert.False(t, list.Items[0].Complete)
		assert.False(t, list.Items[0].SubItems[0].Complete)
		assert.True(t, list.Items[0].SubItems[0].SubItems[0].Complete)
	})

	t.Run("successful deletion of layered list items", func (t *testing.T) {
		deletionPath := fmt.Sprintf("/items/%s", item1Id)
		req, _ := http.NewRequest(http.MethodDelete, deletionPath, nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: item1Id}}
		c.Request = req

		controllers.DeleteItem(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var list models.TodoList
		req, _ = http.NewRequest(http.MethodGet, getListPath, nil)
		w = httptest.NewRecorder()
		c, _ = gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "userId", Value: userId}, {Key: "listId", Value: listId}}
		c.Request = req
		controllers.GetList(c)
		assert.Equal(t, http.StatusOK, w.Code)
		err := json.Unmarshal(w.Body.Bytes(), &list)
		assert.NoError(t, err)
		assert.Empty(t, list.Items)

		getItemPath := fmt.Sprintf("/items/%s", item3Id)
		req, _ = http.NewRequest(http.MethodGet, getItemPath, nil)
		w = httptest.NewRecorder()
		c, _ = gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: item3Id}}
		c.Request = req

		controllers.GetItem(c)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestListSharingDeletions(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db.CleanTestDB()
	db.LoadTestData()

	// share a list with another user - then that user deletes their account
	// list should still be available for other use
	// when the last user deletes their account the list should finally disappear// Assuming userId 1, listId 2, and recipientId 3 exist in your test database
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

	err := services.DeleteUser("tu_1")
	assert.NoError(t, err)

	list, err := services.GetListById("tu_2", "tl_1")
	assert.NoError(t, err)
	t.Log(list)

	tu2Lists, err := services.GetUserListIds("tu_2")
	assert.NoError(t, err)
	t.Log(tu2Lists)

	err = services.DeleteUser("tu_2")
	assert.NoError(t, err)

	var listId string
	row := db.Db.QueryRow("select id from todo_lists where id=$1", "tl_1")
	err = row.Scan(&listId)
	t.Log(listId)
	assert.Error(t, err) // want to make sure list is gone now
}