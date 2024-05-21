package database

func (db *appdbimpl) FollowUser(OwnerID int, FollowedID int) error {
	query := "INSERT INTO FollowDB (OwnerID, FollowedID) VALUES ($1, $2)"

	_, err := db.c.Exec(query, OwnerID, FollowedID)

	if err != nil {
		return err
	}

	return nil
}
