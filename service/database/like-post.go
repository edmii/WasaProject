package database

import "fmt"

func (db *appdbimpl) LikePost(PostID int, OwnerID int) (int, error) {
	var user2ID int

	userQuery := "SELECT OwnerID from PostDB where PostID = $1"
	err := db.c.QueryRow(userQuery, PostID).Scan(&user2ID)
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve user posts: %w", err)
	}
	banExists, err := db.CheckBanStatus(OwnerID, user2ID)

	if err != nil {
		return 0, fmt.Errorf("failed to check ban existence: %w", err)
	}

	if banExists {
		return 0, fmt.Errorf("request failed (user is banned)")
	}

	var exists bool
	checkQuery := `SELECT EXISTS(SELECT 1 FROM LikesDB WHERE LikedPhotoID = $1 AND Ownerid = $2)`
	err = db.c.QueryRow(checkQuery, PostID, OwnerID).Scan(&exists)
	if err != nil {
		return 0, fmt.Errorf("failed to check like existence: %w", err)
	}

	if exists {
		// If the like exists, do nothing
		return 1, nil
	} else {
		// If the like does not exist, insert it
		insertQuery := `INSERT INTO LikesDB (LikedPhotoID, Ownerid) VALUES ($1, $2)`
		_, err := db.c.Exec(insertQuery, PostID, OwnerID)
		if err != nil {
			return 0, fmt.Errorf("failed to insert like: %w", err)
		}
		return 2, nil
	}
}

func (db *appdbimpl) UnlikePost(PostID int, OwnerID int) (int, error) {
	var user2ID int

	userQuery := "SELECT OwnerID from PostDB where PostID = $1"
	err := db.c.QueryRow(userQuery, PostID).Scan(&user2ID)
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve posts owner: %w", err)
	}
	banExists, err := db.CheckBanStatus(OwnerID, user2ID)

	if err != nil {
		return 0, fmt.Errorf("failed to check ban existence: %w", err)
	}

	if banExists {
		return 0, fmt.Errorf("request failed (user is banned)")
	}

	var exists bool
	checkQuery := `SELECT EXISTS(SELECT 1 FROM LikesDB WHERE LikedPhotoID = $1 AND Ownerid = $2)`
	err = db.c.QueryRow(checkQuery, PostID, OwnerID).Scan(&exists)
	if err != nil {
		return 0, fmt.Errorf("failed to check like existence: %w", err)
	}

	if exists {
		// If the like exists, delete it
		deleteQuery := `DELETE FROM LikesDB WHERE LikedPhotoID = $1 AND Ownerid = $2`
		_, err := db.c.Exec(deleteQuery, PostID, OwnerID)
		if err != nil {
			return 0, fmt.Errorf("failed to delete like: %w", err)
		}
		return 1, nil
	} else {
		// If the like does not exist, do nothing
		return 2, nil
	}
}

func (db *appdbimpl) GetLikes(PostID int) ([]int, error) {
	rows, err := db.c.Query("SELECT OwnerID FROM LikesDB WHERE LikedPhotoID = $1", PostID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var likes []int
	for rows.Next() {
		var user int
		if err := rows.Scan(&user); err != nil {
			return nil, err
		}
		likes = append(likes, user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return likes, nil
}

func (db *appdbimpl) GetLikesCount(postID int) (int, error) {
	row := db.c.QueryRow("SELECT COUNT(*) FROM LikesDB WHERE LikedPhotoID = $1", postID)

	// Variable to store the count result
	var count int
	// Scan the result into the count variable
	if err := row.Scan(&count); err != nil {
		return 0, err
	}

	// Return the count
	return count, nil
}
