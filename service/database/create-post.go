package database

func (db *appdbimpl) CreatePost(ownerID int, directory string) error {
	// Prepare the query to insert a new post

	// insert owner id and directory into PostDb

	insertQuery := "INSERT INTO PostDB (OwnerID, Directory) VALUES ($1, $2)"

	// insertQuery2 := "INSERT INTO PostContentDB (PostID, Directory) VALUES ($1, $2)"

	_, err := db.c.Exec(insertQuery, ownerID, directory)

	if err != nil {
		return err
	}

	return nil
}
