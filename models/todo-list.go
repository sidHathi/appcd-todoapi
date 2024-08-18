package models;

type TodoList struct {
	Id string `json:"id"`
	Name string `json:"name"`
	CreatedBy string `json:"CreatedBy"`
	Items []TodoItem `json:"items"`
}

type TodoListSqlRow struct {
	Id string `json:"id"`
	Name string `json:"name"`
	CreatedBy string `json:"CreatedBy"`
}