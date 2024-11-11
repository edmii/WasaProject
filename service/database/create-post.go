package database

import (
	"fmt"
	"time"
)

type Post struct {
	PostID    int    `json:"postID"`
	OwnerID   int    `json:"ownerID"`
	Directory string `json:"imagePath"`
	PostedAt  string `json:"postedAt"`
}

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

func (db *appdbimpl) DeletePost(postID int) error {
	deleteQuery := "DELETE FROM PostDB WHERE PostID = $1"
	_, err := db.c.Exec(deleteQuery, postID)
	if err != nil {
		return fmt.Errorf("failed to delete post: %w", err)
	}
	return nil
}

func (db *appdbimpl) GetUserPosts(username string) ([]Post, error) {

	//("SELECT UserID FROM UserDB WHERE Username = $1", username)
	rows, err := db.c.Query("SELECT * FROM PostDB WHERE OwnerID = (SELECT UserID FROM UserDB WHERE Username = $1)", username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	//return a slice of json, with all info about the post
	var posts []Post
	for rows.Next() {
		var post Post
		if err := rows.Scan(&post.PostID, &post.Directory, &post.OwnerID, &post.PostedAt); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil

}
