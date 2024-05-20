package database

func (db *appdbimpl) CreatePost(ownerID int, directory string) error {

	insertQuery := "INSERT INTO PostDB (OwnerID, Directory) VALUES ($1, $2)"

	_, err := db.c.Exec(insertQuery, ownerID, directory)

	if err != nil {
		return err
	}

	return nil
}
