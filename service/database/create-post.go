package database

import (
	"database/sql"
	"fmt"
	"time"

	structs "github.com/edmii/WasaProject/service/models"
)

// type Post struct {
// 	PostID    int    `json:"postID"`
// 	OwnerID   int    `json:"ownerID"`
// 	Directory string `json:"imagePath"`
// 	PostedAt  string `json:"postedAt"`

// 	RequesterID int `json:"requesterID"`
// }

func (db *appdbimpl) CreatePost(ownerID int, directory string, postedAt time.Time) (int, error) {
	// Step 1: Insert into the PostDB and get the PostID using RETURNING
	insertQuery := "INSERT INTO PostDB (OwnerID, Directory, PostedAt) VALUES ($1, $2, $3) RETURNING PostID"

	var postID int
	err := db.c.QueryRow(insertQuery, ownerID, directory, postedAt).Scan(&postID)
	if err != nil {
		return 0, fmt.Errorf("failed to create post: %w", err)
	}

	// Step 2: Modify the Directory to append the PostID
	updatedDirectory := fmt.Sprintf("%s_%d", directory, postID)

	// Step 3: Update the PostDB record with the new Directory value
	updateQuery := "UPDATE PostDB SET Directory = $1 WHERE PostID = $2"
	_, err = db.c.Exec(updateQuery, updatedDirectory, postID)
	if err != nil {
		return 0, fmt.Errorf("failed to update post directory: %w", err)
	}

	// Return the PostID
	return postID, nil

	// insertQuery := "INSERT INTO PostDB (OwnerID, Directory, PostedAt) VALUES ($1, $2, $3)"

	// _, err := db.c.Exec(insertQuery, ownerID, directory, postedAt)

	// if err != nil {
	// 	return err
	// }

	// return nil
}

func (db *appdbimpl) DeletePost(postID int, requesterID int) error {
	if requesterID <= 0 {
		deleteQuery := "DELETE FROM PostDB WHERE PostID = $1"
		_, err := db.c.Exec(deleteQuery, postID)
		if err != nil {
			return fmt.Errorf("failed to delete post: %w", err)
		}
		return nil
	}

	// First, verify if the requester is the owner of the post
	var ownerID int
	query := "SELECT OwnerID FROM PostDB WHERE PostID = $1"
	err := db.c.QueryRow(query, postID).Scan(&ownerID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("post not found")
		}
		return fmt.Errorf("failed to verify post ownership: %w", err)
	}

	if requesterID != ownerID {
		return fmt.Errorf("unauthorized action: requester is not the owner of the post")
	}

	// Begin a transaction
	tx, err := db.c.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Defer rollback to ensure any error rolls back transaction
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	deletePostQuery := "DELETE FROM PostDB WHERE PostID = $1"
	_, err = tx.Exec(deletePostQuery, postID)
	if err != nil {
		return fmt.Errorf("failed to delete post: %w", err)
	}

	// Query to delete likes related to the post
	deleteLikesQuery := "DELETE FROM LikesDB WHERE PostID = $1"
	_, err = tx.Exec(deleteLikesQuery, postID)
	if err != nil {
		return fmt.Errorf("failed to delete likes: %w", err)
	}

	// Query to delete comments related to the post
	deleteCommentsQuery := "DELETE FROM CommentsDB WHERE PostID = $1"
	_, err = tx.Exec(deleteCommentsQuery, postID)
	if err != nil {
		return fmt.Errorf("failed to delete comments: %w", err)
	}

	// Commit the transaction if all operations succeed
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// deleteQuery := "DELETE FROM PostDB WHERE PostID = $1 AND OwnerID = $2"
// _, err := db.c.Exec(deleteQuery, postID, requesterID)
// if err != nil {
// 	return fmt.Errorf("failed to delete post: %w", err)
// }
// return nil
// }

func (db *appdbimpl) GetUserPosts(username string) ([]structs.Post, error) {

	//("SELECT UserID FROM UserDB WHERE Username = $1", username)
	rows, err := db.c.Query("SELECT * FROM PostDB WHERE OwnerID = (SELECT UserID FROM UserDB WHERE Username = $1) ORDER BY PostedAt DESC", username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	//return a slice of json, with all info about the post
	var posts []structs.Post
	for rows.Next() {
		var post structs.Post
		if err := rows.Scan(&post.PostID, &post.Directory, &post.OwnerID, &post.PostedAt); err != nil {
			return nil, err
		}

		likesCount, err := db.GetLikesCount(post.PostID)
		if err != nil {
			return nil, err
		}
		commentsCount, err := db.GetCommentsCount(post.PostID)
		if err != nil {
			return nil, err
		}
		post.CommentsCount = commentsCount
		post.LikesCount = likesCount
		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil

}

func (db *appdbimpl) GetPostsCount(userID int) (int, error) {
	row := db.c.QueryRow("SELECT COUNT(*) FROM PostDB WHERE OwnerID = $1", userID)

	// Variable to store the count result
	var count int
	// Scan the result into the count variable
	if err := row.Scan(&count); err != nil {
		return 0, err
	}

	// Return the count
	return count, nil
}
