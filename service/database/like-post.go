package database

import "fmt"

func (db *appdbimpl) LikePost(PostID int, OwnerID int) error {
	query := `INSERT INTO LikesDB (LikedPhotoID, Ownerid) VALUES ($1, $2)`

	_, err := db.c.Exec(query, PostID, OwnerID)
	if err != nil {
		return fmt.Errorf("failed to insert like: %w", err)
	}
	return nil
}
