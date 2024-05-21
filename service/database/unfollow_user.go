package database

import "fmt"

func (db *appdbimpl) UnfollowUser(OwnerID int, FollowedID int) error {
	query := "DELETE FROM FollowDB WHERE OwnerID = $1 AND FollowedID= $2"

	_, err := db.c.Exec(query, OwnerID, FollowedID)
	if err != nil {
		return fmt.Errorf("failed to unfollow user: %w", err)
	}
	return nil
}
