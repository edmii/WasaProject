package database

import "fmt"

func (db *appdbimpl) UnbanUser(OwnerID int, PrayID int) error {
	query := "DELETE FROM BanDB WHERE OwnerID = $1 AND PrayID= $2"

	_, err := db.c.Exec(query, OwnerID, PrayID)
	if err != nil {
		return fmt.Errorf("failed to unban user: %w", err)
	}
	return nil
}
