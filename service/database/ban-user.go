package database

func (db *appdbimpl) BanUser(OwnerID int, PrayID int) error {
	query := "INSERT INTO BanDB (OwnerID, PrayID) VALUES ($1, $2)"

	_, err := db.c.Exec(query, OwnerID, PrayID)

	if err != nil {
		return err
	}

	return nil
}
