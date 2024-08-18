package services

import (
	"fmt"
	"todo-api/db"
	"todo-api/models"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

func CreateUser(userData models.CreateUserModel) (*models.User, error) {
	id := uuid.NewString()
	_, err := db.Db.Exec("insert into users(id,name) values ($1,$2);", id, userData.Name)
	if err != nil {
		return nil, err
	}

	return &models.User{ Id: id, Name: userData.Name }, nil
}

func GetUsers() ([]models.User, error) {
	rows, err := db.Db.Query("select (id, name) from users;")
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

func GetUserById(userId string) (*models.User, error) {
	var name string
	row := db.Db.QueryRow("select (name) from users where id=$1;", userId)
	err := row.Scan(&name)
	if err != nil {
		return nil, err
	}

	return &models.User{
		Id: userId,
		Name: name,
	}, nil
}

func GetUserListIds(userId string) ([]string, error) {
	rows, err := db.Db.Query("select (list_id) from user_todo_lists where user_id=$1;", userId)
	if err != nil {
		return nil, err
	}

	listIds := []string{}
	for rows.Next() {
		var listId string
		err = rows.Scan(&listId)
		if err != nil {
			continue
		}
		listIds = append(listIds, listId)
	}
	return listIds, nil
}

func UpdateUser(userId string, updatedData models.CreateUserModel) error {
	_, err := db.Db.Exec("update users set name='$1' where id=$1;", updatedData.Name, userId)
	return err
}

func DeleteUser(userId string) error {
	// may also want to delete any todo lists associated with user?
	// only if the user is the sole owner of said lists
	listIds, err := GetUserListIds(userId)
	if err != nil {
		return err
	}

	for _, lid := range listIds {
		_, err = db.Db.Query("select * from user_todo_lists where list_id<>$1", lid)
		if err != nil {
			// in this case - this is the only user who has access to the list - delete it
			_ = DeleteList(userId, lid)
		}
	}

	_, err = db.Db.Exec("delete from users where id=$1", userId)
	return err
}
