package database

import "fmt"

func (db *appdbimpl) BanUser(OwnerID int, PrayID int) error {
	var exists bool
	checkQuery := `SELECT EXISTS(SELECT 1 FROM BanDB WHERE OwnerID = $1 AND PrayID = $2)`
	err := db.c.QueryRow(checkQuery, OwnerID, PrayID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check ban existence: %w", err)
	}

	if exists {
		// If the ban exists, delete it
		deleteQuery := `DELETE FROM BanDB WHERE OwnerID = $1 AND PrayID = $2`
		_, err := db.c.Exec(deleteQuery, OwnerID, PrayID)
		if err != nil {
			return fmt.Errorf("failed to delete ban: %w", err)
		}
	} else {
		// If the ban does not exist, insert it
		insertQuery := `INSERT INTO BanDB (OwnerID, PrayID) VALUES ($1, $2)`
		_, err := db.c.Exec(insertQuery, OwnerID, PrayID)
		if err != nil {
			return fmt.Errorf("failed to insert ban: %w", err)
		}
	}
	return nil
}
