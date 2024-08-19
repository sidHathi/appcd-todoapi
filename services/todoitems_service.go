package services

import (
	"fmt"
	"todo-api/db"
	"todo-api/models"

	"github.com/google/uuid"
)

func GetItemById(itemId string) (*models.TodoItemNoNest, error) {
	var id string
	var listId string
	var complete bool
	var description string
	var parentId string
	fmt.Println("running get item service with id " + itemId)
	itemRow := db.Db.QueryRow("select id, list_id, description, complete, parent_id from todo_items where id=$1;", itemId)
	err := itemRow.Scan(&id, &listId, &description, &complete, &parentId)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	fmt.Println("Line 24")
	subItemIdRows, err := db.Db.Query("select id from todo_items where parent_id=$1;", itemId)
	if err != nil {
		fmt.Println(err.Error())
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

	fmt.Println("Line 40")
	attachmentRows, err := db.Db.Query("select id, list_id, item_id, s3_url from attachments where item_id=$1", itemId)
	if err != nil {
		fmt.Println(err.Error())
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

	fmt.Println("Line 55")
	fmt.Println(models.TodoItemNoNest{
		Id:          id,
		ListId:      listId,
		Description: description,
		Complete:    complete,
		Attachments: attachments,
		SubItemIds:  sIds,
		ParentId:    parentId,
	})
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

func AddSubItem(item_id string, itemData models.CreateTodoItemData) (*models.TodoItemNoNest, error) {
	var listId string
	parentRow := db.Db.QueryRow("select list_id from todo_items where id=$1;", item_id)
	err := parentRow.Scan(&listId)
	if err != nil {
		return nil, err
	}

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

func UpdateItem(itemId string, updateData models.CreateTodoItemData) error {
	curr, err := GetItemById(itemId)
	if err != nil {
		return err
	}

	newDescription := curr.Description
	if updateData.Description != "" {
		newDescription = updateData.Description
	}
	_, err = db.Db.Exec("update todo_items set description=$1 where id=$2;", newDescription, itemId)
	return err
}

func DeleteItem(itemId string) error {
	item, err := GetItemById(itemId)
	if err != nil {
		return err
	}

	for _, childId := range item.SubItemIds {
		DeleteItem(childId)
	}

	_, err = db.Db.Exec("delete from todo_items where id=$1;", itemId)
	return err
}
