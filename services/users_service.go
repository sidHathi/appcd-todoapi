package services

import (
	"todo-api/db"
	"todo-api/models"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

// Creates a new user
func CreateUser(userData models.CreateUserModel) (*models.User, error) {
	id := uuid.NewString()
	_, err := db.Db.Exec("insert into users(id,name) values ($1,$2);", id, userData.Name)
	if err != nil {
		return nil, err
	}

	return &models.User{Id: id, Name: userData.Name}, nil
}

// Gets all the users
func GetUsers() ([]models.User, error) {
	rows, err := db.Db.Query("select id, name from users;")
	if err != nil {
		return nil, err
	}

	users := []models.User{}
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.Id, &user.Name)
		if err != nil {
			continue
		}
		users = append(users, user)
	}

	return users, nil
}

// Gets a specific user by their id
func GetUserById(userId string) (*models.User, error) {
	var name string
	row := db.Db.QueryRow("select (name) from users where id=$1;", userId)
	err := row.Scan(&name)
	if err != nil {
		return nil, err
	}

	return &models.User{
		Id:   userId,
		Name: name,
	}, nil
}

// Get all the list ids owned be a given user
func GetUserListIds(userId string) ([]string, error) {
	// check that the user actually exists
	_, err := GetUserById(userId)
	if err != nil {
		return nil, err
	}

	// get the list ids
	rows, err := db.Db.Query("select list_id from user_todo_lists where user_id=$1;", userId)
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

// Update a user based on the id
func UpdateUser(userId string, updatedData models.CreateUserModel) error {
	// check that the user exists
	_, err := GetUserById(userId)
	if err != nil {
		return err
	}

	// perform the update
	_, err = db.Db.Exec("update users set name=$1 where id=$2;", updatedData.Name, userId)
	return err
}

// Delete a user with a given id
func DeleteUser(userId string) error {
	// get the lists owned by the user
	listIds, err := GetUserListIds(userId)
	if err != nil {
		return err
	}

	// delete all ownership relations between the user and any lists
	_, err = db.Db.Exec("delete from user_todo_lists where user_id=$1", userId)
	if err != nil {
		return err
	}

	// if the user is the sole remaining owner of the list - delete it
	for _, lid := range listIds {
		var list_owner_id string
		row := db.Db.QueryRow("select user_id from user_todo_lists where list_id=$1 and user_id<>$2", lid, userId)
		err = row.Scan(&list_owner_id)
		if err != nil {
			// in this case - this is the only user who has access to the list - delete it
			_ = DeleteList(userId, lid)
		}
	}

	// delete the user from the user table
	_, err = db.Db.Exec("delete from users where id=$1", userId)
	return err
}
