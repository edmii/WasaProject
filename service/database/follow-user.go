package database

import "fmt"

type Follow struct {
	OwnerID    int `json:"ownerID"`
	FollowedID int `json:"followedID"`
}

func (db *appdbimpl) FollowUser(OwnerID int, FollowedID int) (int, error) {

	var exists bool
	checkQuery := `SELECT EXISTS(SELECT 1 FROM FollowDB WHERE OwnerID = $1 AND FollowedID = $2)`
	err := db.c.QueryRow(checkQuery, OwnerID, FollowedID).Scan(&exists)
	if err != nil {
		return 0, fmt.Errorf("failed to check follow existence: %w", err)
	}

	if exists {
		// If the follow exists, delete it
		deleteQuery := `DELETE FROM FollowDB WHERE OwnerID = $1 AND FollowedID = $2`
		_, err := db.c.Exec(deleteQuery, OwnerID, FollowedID)
		if err != nil {
			return 0, fmt.Errorf("failed to delete follow: %w", err)
		}
		return 1, nil
	} else {
		// If the follow does not exist, insert it
		insertQuery := `INSERT INTO FollowDB (OwnerID, FollowedID) VALUES ($1, $2)`
		_, err := db.c.Exec(insertQuery, OwnerID, FollowedID)
		if err != nil {
			return 0, fmt.Errorf("failed to insert follow: %w", err)
		}
		return 2, nil
	}
}

func (db *appdbimpl) GetFollowers(ownerID int) ([]Follow, error) {
	rows, err := db.c.Query("SELECT ownerID, followedID FROM FollowDB WHERE ownerID = $1", ownerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var followers []Follow
	for rows.Next() {
		var user Follow
		if err := rows.Scan(&user.OwnerID, &user.FollowedID); err != nil {
			return nil, err
		}
		followers = append(followers, user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return followers, nil
}
