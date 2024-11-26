package database

import (
	"fmt"
	"time"

	structs "github.com/edmii/WasaProject/service/models"
)

// type Comment struct {
// 	CommentID int       `json:"commentID"`
// 	Content   string    `json:"content"`
// 	PostID    int       `json:"postID"`
// 	OwnerID   int       `json:"ownerID"`
// 	CreatedAt time.Time `json:"createdAt"`

// 	RequesterID int `json:"requesterID"`
// }

func (db *appdbimpl) CommentPost(PostID int, OwnerID int, Content string, CreatedAt time.Time) error {
	query := "INSERT INTO CommentDB (PhotoID, OwnerID, Content, CreatedAt) VALUES ($1, $2, $3, $4)"
	_, err := db.c.Exec(query, PostID, OwnerID, Content, CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to comment post: %w", err)
	}

	return nil

}

func (db *appdbimpl) DeleteComment(CommentID int, RequesterID int, PostID int) error {
	//execute query only if wonerID is the same as the ownerID of the comment or the ownerID of the post
	query := "DELETE FROM CommentDB WHERE CommentID = $1 AND (OwnerID = $2 OR EXISTS (SELECT 1 FROM PostDB WHERE PostID = $3 AND OwnerID = $2))"

	_, err := db.c.Exec(query, CommentID, RequesterID, PostID)
	if err != nil {
		return fmt.Errorf("failed to delete comment: %w", err)
	}
	return nil
}

func (db *appdbimpl) GetComments(PostID int) ([]structs.Comment, error) {
	rows, err := db.c.Query("SELECT * FROM CommentDB WHERE PhotoID = $1", PostID)
	if err != nil {
		return nil, fmt.Errorf("failed to get comments: %w", err)
	}
	defer rows.Close()

	var comments []structs.Comment
	for rows.Next() {
		var comment structs.Comment
		if err := rows.Scan(&comment.CommentID, &comment.OwnerID, &comment.PostID, &comment.Content, &comment.CreatedAt); err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to get comments: %w", err)
	}

	return comments, nil
}
