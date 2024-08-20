package services

import (
	"todo-api/db"
	"todo-api/models"

	"github.com/google/uuid"
)

// Fetch an item by its id
func GetItemById(itemId string) (*models.TodoItemNoNest, error) {
	// Get the item info
	var id string
	var listId string
	var complete bool
	var description string
	var parentId string
	itemRow := db.Db.QueryRow("select id, list_id, description, complete, parent_id from todo_items where id=$1;", itemId)
	err := itemRow.Scan(&id, &listId, &description, &complete, &parentId)
	if err != nil {
		return nil, err
	}

	// Get the ids of its children
	subItemIdRows, err := db.Db.Query("select id from todo_items where parent_id=$1;", itemId)
	if err != nil {
		return nil, err
	}
	sIds := []string{}
	for subItemIdRows.Next() {
		var sid string
		err := subItemIdRows.Scan(&sid)
		if err != nil {
			continue
		}
		sIds = append(sIds, sid)
	}

	// fetch its attachments
	attachmentRows, err := db.Db.Query("select id, list_id, item_id, s3_url from attachments where item_id=$1", itemId)
	if err != nil {
		return nil, err
	}
	attachments := []models.Attachment{}
	for attachmentRows.Next() {
		var attachment models.Attachment
		err := attachmentRows.Scan(&attachment.Id, &attachment.ListId, &attachment.TodoItemId, &attachment.S3Url)
		if err != nil {
			continue
		}
		attachments = append(attachments, attachment)
	}

	return &models.TodoItemNoNest{
		Id:          id,
		ListId:      listId,
		Description: description,
		Complete:    complete,
		Attachments: attachments,
		SubItemIds:  sIds,
		ParentId:    parentId,
	}, nil
}

// Add a child (sub item) to a todo-list item
func AddSubItem(item_id string, itemData models.CreateTodoItemData) (*models.TodoItemNoNest, error) {
	// fetch the list id of the parent item (also verifies the parent exists)
	var listId string
	parentRow := db.Db.QueryRow("select list_id from todo_items where id=$1;", item_id)
	err := parentRow.Scan(&listId)
	if err != nil {
		return nil, err
	}

	// create and add the new item to the db
	sid := uuid.NewString()
	subItem := models.TodoItemNoNest{
		Id:          sid,
		ListId:      listId,
		Description: itemData.Description,
		Complete:    false,
		Attachments: []models.Attachment{},
		SubItemIds:  []string{},
		ParentId:    item_id,
	}
	_, err = db.Db.Exec("insert into todo_items (id, list_id, description, complete, parent_id) values ($1, $2, $3, $4, $5);", sid, listId, itemData.Description, false, item_id)
	if err != nil {
		return nil, err
	}

	return &subItem, nil
}

// Set an item's completion
func SetItemCompletion(itemId string, completionPayload models.TodoItemCompletionUpdate) (*models.TodoItemNoNest, error) {
	item, err := GetItemById(itemId)
	if err != nil {
		return nil, err
	}

	_, err = db.Db.Exec("update todo_items set complete=$1 where id=$2;", completionPayload.Complete, itemId)
	if err != nil {
		return nil, err
	}

	// if complete - children are complete too
	// if not complete - parent is not complete
	if completionPayload.Complete && len(item.SubItemIds) > 0 {
		for _, childId := range item.SubItemIds {
			SetItemCompletion(childId, completionPayload)
		}
	} else if !completionPayload.Complete && item.ParentId != "" {
		SetItemCompletion(item.ParentId, completionPayload)
	}
	item.Complete = completionPayload.Complete
	return item, nil
}

// Update an item based on its id
func UpdateItem(itemId string, updateData models.CreateTodoItemData) error {
	// check that the item exists
	curr, err := GetItemById(itemId)
	if err != nil {
		return err
	}

	// make sure the update fields are populated
	newDescription := curr.Description
	if updateData.Description != "" {
		newDescription = updateData.Description
	}
	_, err = db.Db.Exec("update todo_items set description=$1 where id=$2;", newDescription, itemId)
	return err
}

// Delete an item
func DeleteItem(itemId string) error {
	// make sure it exists
	item, err := GetItemById(itemId)
	if err != nil {
		return err
	}

	// delete its children
	for _, childId := range item.SubItemIds {
		DeleteItem(childId)
	}

	// delete the item
	_, err = db.Db.Exec("delete from todo_items where id=$1;", itemId)
	return err
}
