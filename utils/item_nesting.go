package utils

import "todo-api/models"

// This file contains utility functions for converting
// todo item data from sql rows into nested todo item lists

// used to nest items within each other - map ensure no duplicate children
type MappedTodoItem struct {
	Id          string                    `json:"id"`
	ListId      string                    `json:"list_id"`
	Description string                    `json:"description"`
	Complete    bool                      `json:"complete"`
	Attachments []models.Attachment       `json:"attachments"`
	SubItems    map[string]MappedTodoItem `json:"sub_items"`
	ParentId    string                    `json:"parent_id"` // empty if no parent
}

// helper function - adds items to a list
func AddItemsToList(list *models.TodoList, items []MappedTodoItem) {
	var parentMap = make(map[string][]MappedTodoItem)
	var idMap = make(map[string]*MappedTodoItem)
	for idx, item := range items {
		subItems, exists := parentMap[item.ParentId]
		if !exists {
			subItems = []MappedTodoItem{}
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

// helper functions - nests a list of items using a relational map
// between the parent items and their children and a relational map
// between the item ids and the actual items
func nestItems(parentMap map[string][]MappedTodoItem, idMap map[string]*MappedTodoItem) []models.TodoItem {
	// want to start with items that have no children -
	// add to their parents and keep going basically
	noChildElems := []*MappedTodoItem{}
	for _, item := range idMap {
		_, exists := parentMap[item.Id]
		if !exists {
			noChildElems = append(noChildElems, item)
		}
	}

	// map ensures no duplication - the last item added should be the one with all the children
	topLevelItems := make(map[string]MappedTodoItem)
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
			// ensure no duplicate subitems=
			parent.SubItems[currItem.Id] = *currItem
			// move parent to back of queue
			itemsToNest = append(itemsToNest, parent)
		} else {
			// if its parent doesn't exist then it's not a subitem
			topLevelItems[currItem.Id] = *currItem
		}
	}

	items := []models.TodoItem{}
	for _, val := range topLevelItems {
		item := ConvertNestingItem(val)
		items = append(items, item)
	}
	return items
}

// converts a nesting NestingTodoItem into a TodoItem - subitems are stored
// in a list rather than a dictionary
func ConvertNestingItem(nestingItem MappedTodoItem) models.TodoItem {
	subItems := []models.TodoItem{}
	for _, val := range nestingItem.SubItems {
		subItems = append(subItems, ConvertNestingItem(val))
	}
	return models.TodoItem{
		Id:          nestingItem.Id,
		ListId:      nestingItem.ListId,
		Description: nestingItem.Description,
		Complete:    nestingItem.Complete,
		Attachments: nestingItem.Attachments,
		SubItems:    subItems,
		ParentId:    nestingItem.ParentId,
	}
}

// Adds a list of attachments to a list of items based on the attachments
// itemId fields
func AddAttachmentsToItems(items []MappedTodoItem, attachments []models.Attachment) []MappedTodoItem {
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
