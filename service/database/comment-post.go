package database

import (
	"fmt"
	"time"
)

type Comment struct {
	ID        int       `json:"id"`
	PostID    int       `json:"postID"`
	OwnerID   int       `json:"ownerID"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created"`
}

func (db *appdbimpl) CommentPost(PostID int, OwnerID int, Content string, CreatedAt time.Time) error {
	query := "INSERT INTO CommentDB (PhotoID, OwnerID, Content, CreatedAt) VALUES ($1, $2, $3, $4)"
	_, err := db.c.Exec(query, PostID, OwnerID, Content, CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to comment post: %w", err)
	}

	return nil

}

func (db *appdbimpl) GetComments(PostID int) ([]Comment, error) {
	rows, err := db.c.Query("SELECT * FROM CommentDB WHERE PhotoID = $1", PostID)
	if err != nil {
		return nil, fmt.Errorf("failed to get comments: %w", err)
	}
	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var comment Comment
		if err := rows.Scan(&comment.ID, &comment.OwnerID, &comment.PostID, &comment.Content, &comment.CreatedAt); err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to get comments: %w", err)
	}

	return comments, nil
}
