package services

import (
	"fmt"
	"todo-api/db"
	"todo-api/models"

	"github.com/google/uuid"
)

func CreateNewList(data models.CreateTodoListData, userId string) (*models.TodoList, error) {
	id := uuid.NewString()
	_, err := db.Db.Exec("insert into todo_lists (id, name, created_by) values ($1, $2, $3);", id, data.Name, userId)
	if err != nil {
		return nil, err
	}

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

func ShareList(userId string, listId string, recipientId string) error {
	_, err := db.Db.Query("select * from user_todo_lists where user_id=$1 and list_id=$2;", userId, listId)
	if err != nil {
		return err
	}

	_, err = db.Db.Exec("insert into user_todo_lists (user_id, list_id) values ($1, $2);", recipientId, listId)
	return err
}

func GetListById(userId string, listId string) (*models.TodoList, error) {
	var id string
	var name string
	var created_by string
	// get the list's info
	row := db.Db.QueryRow("select todo_lists.id, todo_lists.name, todo_lists.created_by from todo_lists inner join user_todo_lists on todo_lists.id = user_todo_lists.list_id where todo_lists.id=$1 and user_todo_lists.user_id=$2;", listId, userId)
	err := row.Scan(&id, &name, &created_by)
	if err != nil {
		fmt.Printf("get list failed at 46 err: %s\n", err.Error())
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
		fmt.Println("get list failed at 60")
		return nil, err
	}
	attachRows, err := db.Db.Query("select id, item_id, s3_url from attachments where list_id=$1;", listId)
	if err != nil {
		fmt.Println("get list failed at 65")
		return nil, err
	}

	var items []models.TodoItem
	for itemRows.Next() {
		var item models.TodoItem
		err := itemRows.Scan(&item.Id, &item.Description, &item.Complete, &item.ParentId)
		if err != nil {
			continue
		}
		item.ListId = listId
		item.Attachments = []models.Attachment{}
		item.SubItems = []models.TodoItem{}
		items = append(items, item)
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

	items = addAttachmentsToItems(items, attachments)
	addItemsToList(&list, items)
	return &list, nil
}

// TODO:
func UpdateListDetails(userId string, listId string, updateData models.CreateTodoListData) error {
	if updateData.Name == "" {
		return nil
	}
	_, err := db.Db.Exec("update todo_lists set name=$1 from user_todo_lists where todo_lists.id=$2 and user_todo_lists.user_id=$3 and user_todo_lists.list_id=$4;", updateData.Name, listId, userId, listId)
	return err
}

func AddListItem(userId string, listId string, item models.CreateTodoItemData) (*models.TodoItem, error) {
	// check to make sure the list belongs to the user
	var testVar string
	row := db.Db.QueryRow("select user_id from user_todo_lists where user_id=$1 and list_id=$2", userId, listId)
	err := row.Scan(&testVar)
	if err != nil {
		return nil, err
	}

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
	// delete all entries in user_todo_lists with the listId
	// delete all attachments with the listId
	// delete all entries in user_todo_items where the item has the listId
	// delete all items with the listid
	// delete the list
	_, err := db.Db.Exec("delete from user_todo_lists where list_id=$1 and user_id=$2;", listId, userId)
	if err != nil {
		return err
	}

	// only want to continue if a list exists for both the given user_id and listId
	var id string
	listRow := db.Db.QueryRow("select user_id from user_todo_lists where list_id=$1;", listId)
	err = listRow.Scan(&id)
	if err == nil {
		// if we're here then another user still has ownership of the list - don't delete it entirely
		return nil
	}

	// _, err = db.Db.Exec("delete from user_todo_items inner join todo_items on user_todo_items.list_id=todo_items.id where todo_items.list_id=$1;", listId)
	// if err != nil {
	// 	return err
	// }

	_, err = db.Db.Exec("delete from attachments where list_id=$1;", listId)
	if err != nil {
		return err
	}

	_, err = db.Db.Exec("delete from todo_items where list_id=$1;", listId)
	if err != nil {
		return err
	}

	_, err = db.Db.Exec("delete from todo_lists where id=$1;", listId)
	if err != nil {
		return err
	}
	return nil
}

func addItemsToList(list *models.TodoList, items []models.TodoItem) {
	var parentMap = make(map[string][]models.TodoItem)
	var idMap = make(map[string]*models.TodoItem)
	for idx, item := range items {
		subItems, exists := parentMap[item.ParentId]
		if !exists {
			subItems = []models.TodoItem{}
		}
		subItems = append(subItems, item)
		parentMap[item.ParentId] = subItems
		idMap[item.Id] = &items[idx]
	}

	nestedItems := nestItems(parentMap, idMap)
	if nestedItems != nil {
		list.Items = nestedItems
	}
}

func nestItems(parentMap map[string][]models.TodoItem, idMap map[string]*models.TodoItem) []models.TodoItem {
	// want to start with items that have no children -
	// add to their parents and keep going basically
	noChildElems := []*models.TodoItem{}
	for _, item := range idMap {
		_, exists := parentMap[item.Id]
		if !exists {
			noChildElems = append(noChildElems, item)
		}
	}

	topLevelItems := make(map[string]models.TodoItem)
	itemsToNest := noChildElems // this functions as a queue
	for len(itemsToNest) > 0 {
		// we guarantee that currItem is populated with all its children
		// because it is either a node without children or the parent
		// of a node whose chlidren have been added to it
		currItem := itemsToNest[0]
		itemsToNest = itemsToNest[1:]
		// add curr item to its parent if it has one
		parentId := currItem.ParentId
		if parentId == "" {
			// if item has no parent then it's not a subitem
			topLevelItems[currItem.Id] = *currItem
			continue
		}

		parent, exists := idMap[parentId]
		if exists {
			parent.SubItems = append(parent.SubItems, *currItem)
			itemsToNest = append(itemsToNest, parent)
		} else {
			// if its parent doesn't exist then it's not a subitem
			topLevelItems[currItem.Id] = *currItem
		}
	}

	items := []models.TodoItem{}
	for _, val := range topLevelItems {
		items = append(items, val)
	}
	return items
}

func addAttachmentsToItems(items []models.TodoItem, attachments []models.Attachment) []models.TodoItem {
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
		if exists && ats != nil {
			items[i].Attachments = ats
		}
	}
	return items
}
