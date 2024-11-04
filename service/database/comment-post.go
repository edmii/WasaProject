package database

import "time"

func (db *appdbimpl) CommentPost(PostID int, OwnerID int, Content string, CreatedAt time.Time) error {
	query := "INSERT INTO CommentDB (PhotoID, OwnerID, Content, CreatedAt) VALUES ($1, $2, $3, $4)"
	_, err := db.c.Exec(query, PostID, OwnerID, Content)

	if err != nil {
		return err
	}

	return nil

}
