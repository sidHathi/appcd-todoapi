package services

import (
	"todo-api/db"
	"todo-api/models"

	"github.com/google/uuid"
)

func CreateNewList(name string, userId string) (*models.TodoList, error) {
	id := uuid.NewString()
	_, err := db.Db.Exec("insert into todo_lists (id, name, created_by) values ($1, '$2', '$3');", id, name, userId)
	if err != nil {
		return nil, err
	}

	return &models.TodoList{
		Id:        id,
		Name:      name,
		CreatedBy: userId,
		Items:     []models.TodoItem{},
	}, nil
}

func GetListById(userId string, listId string) (*models.TodoList, error) {
	var id string
	var name string
	var created_by string
	// get the list's info
	row := db.Db.QueryRow("select (list_id, name, created_by) from todo_lists inner join user_todo_lists on todo_lists.id = user_todo_lists.list_id where todo_lists.id='$1' and user_todo_lists.user_id='$2';", listId, userId)
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
	itemRows, err := db.Db.Query("select (id, description, parent_id) from todo_items where list_id='$1';", listId)
	if err != nil {
		return nil, err
	}
	attachRows, err := db.Db.Query("select (id, item_id, s3_url) from todo_items where list_id='$1';", listId)
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

	items = addAttachmentsToItems(items, attachments)
	addItemsToList(&list, items)
	return &list, nil
}

// TODO:
func UpdateListDetails(userId string, listId string, newName string) error {
	if newName == "" {
		return nil
	}
	_, err := db.Db.Exec("update todo_lists tl set tl.name='$1' from user_todo_lists utl where tl.id = $2 and utl.user_id=$3 and utl.list_id=$4;", newName, listId, userId, listId)
	return err
}

func AddListItem(userId string, listId string, item models.CreateTodoItemData) (*models.TodoItem, error) {
	// check to make sure the list belongs to the user
	var testVar string
	row := db.Db.QueryRow("select * from user_todo_lists where user_id=$1 and list_id=$2", userId, listId)
	err := row.Scan(&testVar)
	if err != nil {
		return nil, err
	}

	id := uuid.NewString()
	completeItem := models.TodoItem{
		Id:          id,
		ListId:      listId,
		Description: item.Description,
		Attachments: []models.Attachment{},
		SubItems:    []models.TodoItem{},
		ParentId:    "",
	}
	_, err = db.Db.Exec("insert into todo_items(id, list_id, description, parent_id) VALUES ($1, $2, $3, $4)", id, listId, item.Description, "")
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
	listRow := db.Db.QueryRow("select * from todo_lists where created_by=$1 and id=$2;", userId, listId)
	err = listRow.Scan(&id)
	if err != nil {
		return err
	}

	_, err = db.Db.Exec("delete from user_todo_items inner join todo_items on user_todo_items.list_id=todo_items.id where todo_items.list_id=$1;", listId)
	if err != nil {
		return err
	}

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
	for _, item := range items {
		subItems, exists := parentMap[item.ParentId]
		if !exists {
			subItems = []models.TodoItem{}
		}
		subItems = append(subItems, item)
		parentMap[item.ParentId] = subItems
		idMap[item.Id] = &item
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
	list.Items = nestItems(parentMap, idMap)
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

	topLevelItems := []models.TodoItem{}
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
			topLevelItems = append(topLevelItems, *currItem)
			continue
		}
		parent, exists := idMap[parentId]
		if exists {
			parent.SubItems = append(parent.SubItems, *currItem)
			itemsToNest = append(itemsToNest, parent)
		} else {
			// if its parent doesn't exist then it's not a subitem
			topLevelItems = append(topLevelItems, *currItem)
		}
	}

	return topLevelItems
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
		if exists {
			items[i].Attachments = ats
		}
	}
	return items
}