package database

import (
	"fmt"
)

func (db *appdbimpl) CreateUser(username string) error {
	// Prepare the query to check if the username already exists
	checkQuery := "SELECT COUNT(*) FROM UserDB WHERE Username = $1"
	var count int
	err := db.c.QueryRow(checkQuery, username).Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		return fmt.Errorf("username already exists")
	}

	// Prepare the query to insert a new user
	insertQuery := "INSERT INTO UserDB (Username) VALUES ($1)"
	_, err = db.c.Exec(insertQuery, username)
	if err != nil {
		return err
	}

	return nil
}
