package database

import (
	"fmt"
)

type Banned struct {
	OwnerID int `json:"ownerID"`
	// list of banned users:
	PrayID int `json:"prayID"`
}

func (db *appdbimpl) BanUser(OwnerID int, PrayID int) (int, error) {
	var exists bool
	checkQuery := `SELECT EXISTS(SELECT 1 FROM BanDB WHERE OwnerID = $1 AND PrayID = $2)`
	err := db.c.QueryRow(checkQuery, OwnerID, PrayID).Scan(&exists)
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

func (db *appdbimpl) GetBannedUsersByOwner(ownerID int) ([]Banned, error) {
	rows, err := db.c.Query("SELECT prayID FROM BanDB WHERE ownerID = $1")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bannedUsers []Banned
	for rows.Next() {
		var user Banned
		if err := rows.Scan(&user.PrayID); err != nil {
			return nil, err
		}
		bannedUsers = append(bannedUsers, user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return bannedUsers, nil
}
