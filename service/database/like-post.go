package database

import "fmt"

func (db *appdbimpl) LikePost(PostID int, OwnerID int) (int, error) {

	var exists bool
	checkQuery := `SELECT EXISTS(SELECT 1 FROM LikesDB WHERE LikedPhotoID = $1 AND Ownerid = $2)`
	err := db.c.QueryRow(checkQuery, PostID, OwnerID).Scan(&exists)
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
		// If the like does not exist, insert it
		insertQuery := `INSERT INTO LikesDB (LikedPhotoID, Ownerid) VALUES ($1, $2)`
		_, err := db.c.Exec(insertQuery, PostID, OwnerID)
		if err != nil {
			return 0, fmt.Errorf("failed to insert like: %w", err)
		}
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
