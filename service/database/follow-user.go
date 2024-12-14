package database

import "fmt"

// type Follow struct {
// 	OwnerID    int `json:"ownerID"`
// 	FollowedID int `json:"followedID"`
// }

func (db *appdbimpl) FollowUser(OwnerID int, FollowedID int) (int, error) {

	var exists bool
	var ownerExists, followedExists bool

	banExists, err := db.CheckBanStatus(OwnerID, FollowedID)

	if err != nil {
		return 0, fmt.Errorf("failed to check ban existence: %w", err)
	}

	if banExists {
		return 0, fmt.Errorf("request failed (user is banned)")
	}

	// Check if Owner exists in UserDB
	checkOwnerQuery := `SELECT EXISTS(SELECT 1 FROM UserDB WHERE UserID = $1)`
	err = db.c.QueryRow(checkOwnerQuery, OwnerID).Scan(&ownerExists)
	if err != nil {
		return 0, fmt.Errorf("failed to check owner existence: %w", err)
	}

	// Check if Followed exists in UserDB
	checkPrayQuery := `SELECT EXISTS(SELECT 1 FROM UserDB WHERE UserID = $2)`
	err = db.c.QueryRow(checkPrayQuery, FollowedID).Scan(&followedExists)
	if err != nil {
		return 0, fmt.Errorf("failed to check followed existence: %w", err)
	}

	if !ownerExists || !followedExists {
		return 0, fmt.Errorf("one or both users do not exist")
	}

	checkQuery := `SELECT EXISTS(SELECT 1 FROM FollowDB WHERE OwnerID = $1 AND FollowedID = $2)`
	err = db.c.QueryRow(checkQuery, OwnerID, FollowedID).Scan(&exists)
	if err != nil {
		return 0, fmt.Errorf("failed to check follow existence: %w", err)
	}

	if exists {
		// If the follow exists, do nothing
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

func (db *appdbimpl) UnfollowUser(OwnerID int, FollowedID int) (int, error) {

	var exists bool
	var ownerExists, followedExists bool

	banExists, err := db.CheckBanStatus(OwnerID, FollowedID)

	if err != nil {
		return 0, fmt.Errorf("failed to check ban existence: %w", err)
	}

	if banExists {
		return 0, fmt.Errorf("request failed (user is banned)")
	}

	// Check if Owner exists in UserDB
	checkOwnerQuery := `SELECT EXISTS(SELECT 1 FROM UserDB WHERE UserID = $1)`
	err = db.c.QueryRow(checkOwnerQuery, OwnerID).Scan(&ownerExists)
	if err != nil {
		return 0, fmt.Errorf("failed to check owner existence: %w", err)
	}

	// Check if Followed exists in UserDB
	checkPrayQuery := `SELECT EXISTS(SELECT 1 FROM UserDB WHERE UserID = $2)`
	err = db.c.QueryRow(checkPrayQuery, FollowedID).Scan(&followedExists)
	if err != nil {
		return 0, fmt.Errorf("failed to check followed existence: %w", err)
	}

	if !ownerExists || !followedExists {
		return 0, fmt.Errorf("one or both users do not exist")
	}

	checkQuery := `SELECT EXISTS(SELECT 1 FROM FollowDB WHERE OwnerID = $1 AND FollowedID = $2)`
	err = db.c.QueryRow(checkQuery, OwnerID, FollowedID).Scan(&exists)
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
		// If the follow does not exist, do nothing
		return 2, nil
	}
}

func (db *appdbimpl) GetFollowed(ownerID int) ([]int, error) {
	rows, err := db.c.Query("SELECT followedID FROM FollowDB WHERE ownerID = $1", ownerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var followers []int
	for rows.Next() {
		var user int
		if err := rows.Scan(&user); err != nil {
			return nil, err
		}
		followers = append(followers, user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return followers, nil
}

func (db *appdbimpl) GetFollowers(ownerID int) ([]int, error) {
	rows, err := db.c.Query("SELECT ownerID FROM FollowDB WHERE followedID = $1", ownerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var followers []int
	for rows.Next() {
		var user int
		if err := rows.Scan(&user); err != nil {
			return nil, err
		}
		followers = append(followers, user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return followers, nil
}

func (db *appdbimpl) GetFollowersCount(ownerID int) (int, error) {
	row := db.c.QueryRow("SELECT COUNT(*) FROM FollowDB WHERE FollowedID = $1", ownerID)

	// Variable to store the count result
	var count int
	// Scan the result into the count variable
	if err := row.Scan(&count); err != nil {
		return 0, err
	}

	// Return the count
	return count, nil
}

func (db *appdbimpl) GetFollowedCount(ownerID int) (int, error) {
	row := db.c.QueryRow("SELECT COUNT(*) FROM FollowDB WHERE OwnerID = $1", ownerID)

	// Variable to store the count result
	var count int
	// Scan the result into the count variable
	if err := row.Scan(&count); err != nil {
		return 0, err
	}

	// Return the count
	return count, nil
}
