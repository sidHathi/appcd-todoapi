package services

import (
	"todo-api/db"
	"todo-api/models"

	"github.com/google/uuid"
)

func CreateNewList(name string, creator string) (*models.TodoList, error) {
	id := uuid.NewString()
	_, err := db.Db.Exec("insert into todo_lists (id, name, created_by) values ($1, $2, $3);", id, name, creator)
	if err != nil {
		return nil, err
	}

	return &models.TodoList{
		Id: id,
		Name: name,
		CreatedBy: creator,
		Items: []models.TodoItem {},
	}, nil
}

func GetListById(userId string, listId string) (*models.TodoList, error) {
	var id string
	var name string
	var created_by string
	// get the list's info
	row := db.Db.QueryRow("select (list_id, name, created_by) from todo_lists left join user_todo_lists on todo_lists.id = user_todo_lists.list_id where todo_lists.id=$1 and user_todo_lists.user_id=$2;", listId, userId)
	err := row.Scan(&id, &name, &created_by)
	if err != nil {
		return nil, err
	}
	list := models.TodoList {
		Id: id,
		Name: name,
		CreatedBy: created_by,
		Items: []models.TodoItem{},
	}

	// get the list's items
	itemRows, err := db.Db.Query("select (id, description, parent_id) from todo_items where list_id=$1", listId)
	if err != nil {
		return nil, err
	}
	attachRows, err := db.Db.Query("select (id, item_id, s3_url) from todo_items where list_id=$1", listId)
	if err != nil {
		return nil, err
	}

	var items []models.TodoItem
	for itemRows.Next() {
		var item models.TodoItem
		err := itemRows.Scan(&item.Id, &item.Description, &item.ParentId)
		if err != nil {
			continue
		}
		items = append(items, item)
	}

	var attachments []models.Attachment
	for attachRows.Next() {
		var attachment models.Attachment
		err := attachRows.Scan(&attachment.Id, &attachment.TodoItemId, &attachment.S3Url)
		if err != nil {
			continue
		}
		attachments = append(attachments, attachment)
	}

	items = AddAttachmentsToItems(items, attachments)
	AddItemsToList(&list, items)
	return &list, nil
}

func AddItemsToList(list *models.TodoList, items []models.TodoItem) {
	var parentMap = make(map[string][]models.TodoItem)
	var idMap = make(map[string]models.TodoItem)
	for _, item := range items {
		subItems, exists := parentMap[item.ParentId]
		if !exists {
			subItems = []models.TodoItem{}
		}
		subItems = append(subItems, item)
		parentMap[item.ParentId] = subItems
		idMap[item.Id] = item
	}

	for parentId, item := range parentMap {
		if parentId == "" {
			continue
		}
		parent, exists := idMap[parentId]
		if !exists {
			continue
		}

		parent.SubItems = append(parent.SubItems, item...)
		idMap[parentId] = parent
	}

	var topLevelItems []models.TodoItem
	for _, item := range idMap {
		if item.ParentId == "" {
			topLevelItems = append(topLevelItems, item)
		}
	}
	list.Items = topLevelItems
}

func AddAttachmentsToItems(items []models.TodoItem, attachments []models.Attachment) []models.TodoItem {
	atItemMap := make(map[string][]models.Attachment)
	for _, at := range attachments {
		ats, exists := atItemMap[at.TodoItemId]
		if !exists {
			ats = []models.Attachment{}
		}
		ats = append(ats, at)
		atItemMap[at.TodoItemId] = ats
	}

	for i := range items {
		ats, exists := atItemMap[items[i].Id]
		if exists {
			items[i].Attachments = ats
		}
	}
	return items
}
