package models;

// User represents a user in the system.
// swagger:model
type User struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

// CreateUserModel represents the input data needed to create a new user.
// swagger:model
type CreateUserModel struct {
	Name string `json:"name"`
}