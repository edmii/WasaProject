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
	var user2ID int

	userQuery := "SELECT OwnerID from PostDB where PostID $1"
	err := db.c.QueryRow(userQuery, PostID).Scan(&user2ID)
	banExists, err := db.CheckBanStatus(OwnerID, user2ID)

	if err != nil {
		return 0, fmt.Errorf("failed to check ban existence: %w", err)
	}

	if banExists {
		return 0, fmt.Errorf("request failed (user is banned)")
	}

	query := "INSERT INTO CommentDB (PhotoID, OwnerID, Content, CreatedAt) VALUES ($1, $2, $3, $4)"
	_, err = db.c.Exec(query, PostID, OwnerID, Content, CreatedAt)

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

func (db *appdbimpl) GetComments(postID int) ([]structs.Comment, error) {
	rows, err := db.c.Query("SELECT * FROM CommentDB WHERE PhotoID = $1", postID)
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

func (db *appdbimpl) GetCommentsCount(postID int) (int, error) {
	row := db.c.QueryRow("SELECT COUNT(*) FROM CommentDB WHERE PhotoID = $1", postID)

	// Variable to store the count result
	var count int
	// Scan the result into the count variable
	if err := row.Scan(&count); err != nil {
		return 0, err
	}

	// Return the count
	return count, nil
}
