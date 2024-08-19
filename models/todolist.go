package models

// TodoList represents a user's todo list in the system.
// swagger:model
type TodoList struct {
	Id        string     `json:"id"`
	Name      string     `json:"name"`
	CreatedBy string     `json:"CreatedBy"`
	Items     []TodoItem `json:"items"`
}

// CreateTodoListData represents the data needed to create a todo list
// swagger:model
type CreateTodoListData struct {
	Name string `json:"name"`
}

// internal representation of the way todo lists are stored in postgres
type TodoListSqlRow struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	CreatedBy string `json:"CreatedBy"`
}

// ShareListData represents the data needed to share one user's todo
// list with another user
// swagger:model
type ShareListData struct {
	RecipientId string `json:"recipientId"`
}
