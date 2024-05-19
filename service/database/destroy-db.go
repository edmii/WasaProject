package database

import (
	"fmt"
)

func (db *appdbimpl) DestroyDB() error {
	// Prepare the query to drop all tables
	dropTable1 := "DROP TABLE IF EXISTS UserDB"
	_, err := db.c.Exec(dropTable1)
	if err != nil {
		return fmt.Errorf("error dropping table: %w", err)
	}

	dropTable2 := "DROP TABLE IF EXISTS BanDB"
	_, err = db.c.Exec(dropTable2)
	if err != nil {
		return fmt.Errorf("error dropping table: %w", err)
	}

	dropTable3 := "DROP TABLE IF EXISTS PrayDB"
	_, err = db.c.Exec(dropTable3)
	if err != nil {
		return fmt.Errorf("error dropping table: %w", err)
	}

	dropTable4 := "DROP TABLE IF EXISTS LikesDB"
	_, err = db.c.Exec(dropTable4)
	if err != nil {
		return fmt.Errorf("error dropping table: %w", err)
	}

	dropTable5 := "DROP TABLE IF EXISTS CommentDB"
	_, err = db.c.Exec(dropTable5)
	if err != nil {
		return fmt.Errorf("error dropping table: %w", err)
	}

	dropTable6 := "DROP TABLE IF EXISTS PostDB"
	_, err = db.c.Exec(dropTable6)
	if err != nil {
		return fmt.Errorf("error dropping table: %w", err)
	}

	return nil
}
