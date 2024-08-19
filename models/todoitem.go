package models

// TodoItem represents a todo list item in the system.
// swagger:model
type TodoItem struct {
	Id          string       `json:"id"`
	ListId      string       `json:"list_id"`
	Description string       `json:"description"`
	Complete    bool         `json:"complete"`
	Attachments []Attachment `json:"attachments"`
	SubItems    []TodoItem   `json:"sub_items"`
	ParentId    string       `json:"parent_id"` // empty if no parent
}

// TodoItemNoNest represents an unnested todo list item in the system.
// Instead of including all of it's children, it only contains
// their ids referenced as strings. Returned for efficiency reasons
// swagger:model
type TodoItemNoNest struct {
	Id          string       `json:"id"`
	ListId      string       `json:"list_id"`
	Description string       `json:"description"`
	Complete    bool         `json:"complete"`
	Attachments []Attachment `json:"attachments"`
	SubItemIds  []string     `json:"sub_item_ids"`
	ParentId    string       `json:"parent_id"` // empty if no parent
}

// CreateTodoItemData represents the input data used to create a new todo item
// swagger:model
type CreateTodoItemData struct {
	Description string `json:"description"`
}

// TodoItemCompletionUpdate represents the input data used to
// update the completion status of a todo item
// swagger:model
type TodoItemCompletionUpdate struct {
	Complete bool `json:"complete"`
}
