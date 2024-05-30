package database

import "fmt"

func (db *appdbimpl) LikePost(PostID int, OwnerID int) error {

	var exists bool
	checkQuery := `SELECT EXISTS(SELECT 1 FROM LikesDB WHERE LikedPhotoID = $1 AND Ownerid = $2)`
	err := db.c.QueryRow(checkQuery, PostID, OwnerID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check like existence: %w", err)
	}

	if exists {
		// If the like exists, delete it
		deleteQuery := `DELETE FROM LikesDB WHERE LikedPhotoID = $1 AND Ownerid = $2`
		_, err := db.c.Exec(deleteQuery, PostID, OwnerID)
		if err != nil {
			return fmt.Errorf("failed to delete like: %w", err)
		}
	} else {
		// If the like does not exist, insert it
		insertQuery := `INSERT INTO LikesDB (LikedPhotoID, Ownerid) VALUES ($1, $2)`
		_, err := db.c.Exec(insertQuery, PostID, OwnerID)
		if err != nil {
			return fmt.Errorf("failed to insert like: %w", err)
		}
	}
	return nil
}
