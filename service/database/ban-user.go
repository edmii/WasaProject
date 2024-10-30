package database

import (
	"fmt"
)

type Banned struct {
	OwnerID int `json:"ownerID"`
	PrayID  int `json:"prayID"`
}

func (db *appdbimpl) BanUser(OwnerID int, PrayID int) (int, error) {
	var exists bool

	// check if owner and pray exist in db
	checkQuery1 := `SELECT EXISTS(SELECT 1 FROM UserDB WHERE UserID = $1) AND EXISTS(SELECT 1 FROM UserDB WHERE UserID = $2)`
	err := db.c.QueryRow(checkQuery1, OwnerID, PrayID).Scan(&exists)
	if err != nil {
		return 0, fmt.Errorf("failed to check user existence: %w", err)
	}

	if !exists {
		return 0, fmt.Errorf("user does not exist")
	}

	checkQuery := `SELECT EXISTS(SELECT 1 FROM BanDB WHERE OwnerID = $1 AND PrayID = $2)`
	err = db.c.QueryRow(checkQuery, OwnerID, PrayID).Scan(&exists)
	if err != nil {
		return 0, fmt.Errorf("failed to check ban existence: %w", err)
	}

	if exists {
		// If the ban exists, delete it
		deleteQuery := `DELETE FROM BanDB WHERE OwnerID = $1 AND PrayID = $2`
		_, err := db.c.Exec(deleteQuery, OwnerID, PrayID)
		if err != nil {
			return 0, fmt.Errorf("failed to delete ban: %w", err)
		}
		return 1, nil
	} else {
		// If the ban does not exist, insert it
		insertQuery := `INSERT INTO BanDB (OwnerID, PrayID) VALUES ($1, $2)`
		_, err := db.c.Exec(insertQuery, OwnerID, PrayID)
		if err != nil {
			return 0, fmt.Errorf("failed to insert ban: %w", err)
		}
		return 2, nil
	}

}

func (db *appdbimpl) GetBannedUsers(ownerID int) ([]int, error) {
	rows, err := db.c.Query("SELECT prayID FROM BanDB WHERE ownerID = $1", ownerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bannedUsers []int
	for rows.Next() {
		var PrayID int
		if err := rows.Scan(&PrayID); err != nil {
			return nil, err
		}
		bannedUsers = append(bannedUsers, PrayID)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return bannedUsers, nil
}
