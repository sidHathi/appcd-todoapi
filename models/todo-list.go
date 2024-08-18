package models

type TodoList struct {
	Id        string     `json:"id"`
	Name      string     `json:"name"`
	CreatedBy string     `json:"CreatedBy"`
	Items     []TodoItem `json:"items"`
}

type CreateTodoListData struct {
	Name string `json:"name"`
}

type TodoListSqlRow struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	CreatedBy string `json:"CreatedBy"`
}
