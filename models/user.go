package models;

type User struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type CreateUserModel struct {
	Name string `json:"name"`
}