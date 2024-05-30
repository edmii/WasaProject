package database

import "fmt"

func (db *appdbimpl) FollowUser(OwnerID int, FollowedID int) error {

	var exists bool
	checkQuery := `SELECT EXISTS(SELECT 1 FROM FollowDB WHERE OwnerID = $1 AND FollowedID = $2)`
	err := db.c.QueryRow(checkQuery, OwnerID, FollowedID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check follow existence: %w", err)
	}

	if exists {
		// If the follow exists, delete it
		deleteQuery := `DELETE FROM FollowDB WHERE OwnerID = $1 AND FollowedID = $2`
		_, err := db.c.Exec(deleteQuery, OwnerID, FollowedID)
		if err != nil {
			return fmt.Errorf("failed to delete follow: %w", err)
		}
	} else {
		// If the follow does not exist, insert it
		insertQuery := `INSERT INTO FollowDB (OwnerID, FollowedID) VALUES ($1, $2)`
		_, err := db.c.Exec(insertQuery, OwnerID, FollowedID)
		if err != nil {
			return fmt.Errorf("failed to insert follow: %w", err)
		}
	}
	return nil
}
