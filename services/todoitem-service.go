package services

import (
	"todo-api/db"
	"todo-api/models"

	"github.com/google/uuid"
)

func GetItemById(item_id string) (*models.TodoItemNoNest, error) {
	var id string
	var listId string
	var description string
	var parentId string
	itemRow := db.Db.QueryRow("select (id, list_id, description, parent_id) from todo_items where id=$1;", item_id)
	err := itemRow.Scan(&id, &listId, &description, &parentId)
	if err != nil {
		return nil, err
	}

	subItemIdRows, err := db.Db.Query("select (id) from todo_items where parent_id=$1;", item_id)
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

	attachmentRows, err := db.Db.Query("select (id, list_id, item_id, s3_url) from attachments where todo_item_id=$1", item_id)
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
		Attachments: attachments,
		SubItemIds:  sIds,
		ParentId:    parentId,
	}, nil
}

func AddSubItem(item_id string, itemData models.CreateTodoItemData) (*models.TodoItemNoNest, error) {
	var listId string
	parentRow := db.Db.QueryRow("select (list_id) from todo_items where id=$1;", item_id)
	err := parentRow.Scan(&listId)
	if err != nil {
		return nil, err
	}

	sid := uuid.NewString()
	subItem := models.TodoItemNoNest {
		Id: sid,
		ListId: listId,
		Description: itemData.Description,
		Attachments: []models.Attachment{},
		SubItemIds: []string{},
		ParentId: item_id,
	}
	_, err = db.Db.Exec("insert into todo_items (id, list_id, description, parent_id) values ($1, $2, '$3', $4);", sid, listId, itemData.Description, item_id)
	if err != nil {
		return nil, err
	}

	return &subItem, nil
}

func UpdateItem(itemId string, updateData models.CreateTodoItemData) error {
	_, err := db.Db.Exec("update todo_items set description='$1' where id=$2", updateData.Description, itemId)
	return err
}

func DeleteItem(itemId string) error {
	_, err := db.Db.Exec("delete todo_items where id=$1", itemId)
	return err
}
