package database

func (db *appdbimpl) CommentPost(PostID int, OwnerID int, Content string) error {
	query := "INSERT INTO CommentDB (PhotoID, OwnerID, Content) VALUES ($1, $2, $3)"

	_, err := db.c.Exec(query, PostID, OwnerID, Content)

	if err != nil {
		return err
	}

	return nil

}
