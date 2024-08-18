package models

type TodoItem struct {
	Id string `json:"id"`
	ListId string `json:"list_id"`
	Description string `json:"description"`
	Attachments []Attachment `json:"attachments"`
	SubItems []TodoItem `json:"sub_items"`
	ParentId string `json:"parent_id"` // empty if no parent
}

type TodoItemNoNest struct {
	Id string `json:"id"`
	ListId string `json:"list_id"`
	Description string `json:"description"`
	Attachments []Attachment `json:"attachments"`
	SubItemIds []string `json:"sub_items"`
	ParentId string `json:"parent_id"` // empty if no parent
}

type CreateTodoItemData struct {
	Description string `json:"description"`
}