package database

import (
	"fmt"
)

func (db *appdbimpl) CheckUserExist(username string) (bool, error) {
	var exist bool
	// Prepare the query to check if the username already exists
	checkQuery := "SELECT EXIST(SELECT * FROM UserDB WHERE Username = $1)"
	err := db.c.QueryRow(checkQuery, username).Scan(&exist)
	if err != nil {
		return false, err
	}

	return exist, nil
}

func (db *appdbimpl) CreateUser(username string) error {
	// Prepare the query to insert a new user
	insertQuery := "INSERT INTO UserDB (Username) VALUES ($1)"
	_, err := db.c.Exec(insertQuery, username)
	if err != nil {
		return err
	}

	return nil
}

func (db *appdbimpl) GetUserID(username string) (int, error) {
	// Prepare the query to get the user ID
	selectQuery := "SELECT UserID FROM UserDB WHERE Username = $1"
	var userID int
	err := db.c.QueryRow(selectQuery, username).Scan(&userID)
	if err != nil {
		return 0, err
	}

	return userID, nil
}

func (db *appdbimpl) GetUsername(userID int) string {
	// Prepare the query to get the username
	selectQuery := "SELECT Username FROM UserDB WHERE UserID = $1"
	var username string
	err := db.c.QueryRow(selectQuery, userID).Scan(&username)
	if err != nil {
		return ""
	}

	return username
}

func (db *appdbimpl) ChangeUsername(userID int, newUsername string) error {
	updateQuery := "UPDATE UserDB SET Username = $1 WHERE UserID = $2"

	// Execute the update query with the provided new username and user ID
	_, err := db.c.Exec(updateQuery, newUsername, userID)
	if err != nil {
		return fmt.Errorf("failed to update username: %w", err)
	}

	return nil
}
