package services

import (
	"fmt"
	"todo-api/db"
	"todo-api/models"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

func CreateUser(name string) (*models.User, error) {
	id := uuid.NewString()
	_, err := db.Db.Exec("insert into users(id,name) values ($1,$2);", id, name)
	if err != nil {
		return nil, err
	}

	return &models.User{ Id: id, Name: name }, nil
}

func GetUsers() ([]models.User, error) {
	rows, err := db.Db.Query("select * from users;")
	if err != nil {
		return nil, err
	}

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.Id, &user.Name)
		if err != nil {
			fmt.Println("Error scanning row:", err)
			continue
		}
		users = append(users, user)
	}

	return users, nil
}
