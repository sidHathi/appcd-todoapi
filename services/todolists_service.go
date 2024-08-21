package services

import (
	"todo-api/db"
	"todo-api/models"
	"todo-api/utils"

	"github.com/google/uuid"
)

// Create a new list for a given user
func CreateNewList(data models.CreateTodoListData, userId string) (*models.TodoList, error) {
	// check to make sure the user exists
	_, err := GetUserById(userId)
	if err != nil {
		return nil, err
	}

	// add the list to the db
	id := uuid.NewString()
	_, err = db.Db.Exec("insert into todo_lists (id, name, created_by) values ($1, $2, $3);", id, data.Name, userId)
	if err != nil {
		return nil, err
	}

	// add the ownership relation between the user and the list to the db
	_, err = db.Db.Exec("insert into user_todo_lists (user_id, list_id) values ($1, $2);", userId, id)
	if err != nil {
		return nil, err
	}

	return &models.TodoList{
		Id:        id,
		Name:      data.Name,
		CreatedBy: userId,
		Items:     []models.TodoItem{},
	}, nil
}

// Share a list owned by one user with another (recipient) user
func ShareList(userId string, listId string, recipientId string) error {
	// make sure that the user actually owns the input list
	var _testExists string
	err := db.Db.QueryRow("select user_id from user_todo_lists where user_id=$1 and list_id=$2;", userId, listId).Scan(&_testExists)
	if err != nil {
		return err
	}

	// add the new ownership relation to the db
	_, err = db.Db.Exec("insert into user_todo_lists (user_id, list_id) values ($1, $2);", recipientId, listId)
	return err
}

// Get a list (and all nested items) by its id
func GetListById(userId string, listId string) (*models.TodoList, error) {
	// get the list's info
	var id string
	var name string
	var created_by string
	row := db.Db.QueryRow("select todo_lists.id, todo_lists.name, todo_lists.created_by from todo_lists inner join user_todo_lists on todo_lists.id = user_todo_lists.list_id where todo_lists.id=$1 and user_todo_lists.user_id=$2;", listId, userId)
	err := row.Scan(&id, &name, &created_by)
	if err != nil {
		return nil, err
	}
	list := models.TodoList{
		Id:        id,
		Name:      name,
		CreatedBy: created_by,
		Items:     []models.TodoItem{},
	}

	// get the list's items
	itemRows, err := db.Db.Query("select id, description, complete, parent_id from todo_items where list_id=$1;", listId)
	if err != nil {
		return nil, err
	}
	var items []utils.MappedTodoItem
	for itemRows.Next() {
		var item utils.MappedTodoItem
		err := itemRows.Scan(&item.Id, &item.Description, &item.Complete, &item.ParentId)
		if err != nil {
			continue
		}
		item.ListId = listId
		item.Attachments = []models.Attachment{}
		item.SubItems = make(map[string]utils.MappedTodoItem)
		items = append(items, item)
	}

	// get the list's attachments
	attachRows, err := db.Db.Query("select id, item_id, s3_url from attachments where list_id=$1;", listId)
	if err != nil {
		return nil, err
	}
	var attachments []models.Attachment
	for attachRows.Next() {
		var attachment models.Attachment
		err := attachRows.Scan(&attachment.Id, &attachment.TodoItemId, &attachment.S3Url)
		if err != nil {
			continue
		}
		attachment.ListId = listId
		attachments = append(attachments, attachment)
	}

	// add the attachments to their parent items
	items = utils.AddAttachmentsToItems(items, attachments)
	// nest the items within their parents and add the top-level ones to the list
	utils.AddItemsToList(&list, items)
	return &list, nil
}

// Update a list's info
func UpdateListDetails(userId string, listId string, updateData models.CreateTodoListData) error {
	// check to make sure the list exists and is owned by the given user
	var _testExists string
	err := db.Db.QueryRow("select id from todo_lists where id=$1", listId).Scan(&_testExists)
	if err != nil {
		return err
	}
	err = db.Db.QueryRow("select user_id from user_todo_lists where list_id=$1 and user_id=$2", listId, userId).Scan(&_testExists)
	if err != nil {
		return err
	}

	// Update the item based on its name
	if updateData.Name == "" {
		return nil
	}
	_, err = db.Db.Exec("update todo_lists set name=$1 from user_todo_lists where todo_lists.id=$2 and user_todo_lists.user_id=$3 and user_todo_lists.list_id=$4;", updateData.Name, listId, userId, listId)
	return err
}

// Add an item to a list
func AddListItem(userId string, listId string, item models.CreateTodoItemData) (*models.TodoItem, error) {
	// check to make sure the list exists and belongs to the user
	var testVar string
	err := db.Db.QueryRow("select user_id from user_todo_lists where user_id=$1 and list_id=$2", userId, listId).Scan(&testVar)
	if err != nil {
		return nil, err
	}

	// create and add the item
	id := uuid.NewString()
	completeItem := models.TodoItem{
		Id:          id,
		ListId:      listId,
		Description: item.Description,
		Complete:    false,
		Attachments: []models.Attachment{},
		SubItems:    []models.TodoItem{},
		ParentId:    "",
	}
	_, err = db.Db.Exec("insert into todo_items (id, list_id, description, complete, parent_id) VALUES ($1, $2, $3, $4, $5)", id, listId, item.Description, false, "")
	if err != nil {
		return nil, err
	}
	return &completeItem, nil
}

func DeleteList(userId string, listId string) error {
	// remove the user's ownership of the list
	_, err := db.Db.Exec("delete from user_todo_lists where list_id=$1 and user_id=$2;", listId, userId)
	if err != nil {
		return err
	}

	// if another user also own's the list, they should still have it - do not proceed with deletion
	var id string
	listRow := db.Db.QueryRow("select user_id from user_todo_lists where list_id=$1;", listId)
	err = listRow.Scan(&id)
	if err == nil {
		// if we're here then another user still has ownership of the list - don't delete it entirely
		return nil
	}
	// if we make it past here - the list should be removed entirely

	// delete the list's attachments
	_, err = db.Db.Exec("delete from attachments where list_id=$1;", listId)
	if err != nil {
		return err
	}

	// delete the list's items
	_, err = db.Db.Exec("delete from todo_items where list_id=$1;", listId)
	if err != nil {
		return err
	}

	// delete the list itself
	_, err = db.Db.Exec("delete from todo_lists where id=$1;", listId)
	if err != nil {
		return err
	}
	return nil
}
