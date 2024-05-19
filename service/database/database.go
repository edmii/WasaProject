/*
Package database is the middleware between the app database and the code. All data (de)serialization (save/load) from a
persistent database are handled here. Database specific logic should never escape this package.

To use this package you need to apply migrations to the database if needed/wanted, connect to it (using the database
data source name from config), and then initialize an instance of AppDatabase from the DB connection.

For example, this code adds a parameter in `webapi` executable for the database data source name (add it to the
main.WebAPIConfiguration structure):

	DB struct {
		Filename string `conf:""`
	}

This is an example on how to migrate the DB and connect to it:

	// Start Database
	logger.Println("initializing database support")
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		logger.WithError(err).Error("error opening SQLite DB")
		return fmt.Errorf("opening SQLite: %w", err)
	}
	defer func() {
		logger.Debug("database stopping")
		_ = db.Close()
	}()

Then you can initialize the AppDatabase and pass it to the api package.
*/
package database

import (
	"database/sql"
	"errors"
	"fmt"
)

// AppDatabase is the high level interface for the DB
type AppDatabase interface {
	GetName() (string, error)
	SetName(name string) error
	GetDatabaseTableContent(tableName string) ([]map[string]interface{}, error)

	Ping() error
}

type appdbimpl struct {
	c *sql.DB
}

// New returns a new instance of AppDatabase based on the SQLite connection `db`.
// `db` is required - an error will be returned if `db` is `nil`.
func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building a AppDatabase")
	}

	// Check if table exists. If not, the database is empty, and we need to create the structure
	var tableName string
	err := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='example_table';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		sqlStmt := `CREATE TABLE example_table (id INTEGER NOT NULL PRIMARY KEY, name TEXT);`

		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}

		// Create the User table
		createTableSQL := `
			CREATE TABLE IF NOT EXISTS UserDB (
				UserID INT AUTO_INCREMENT,
				Username VARCHAR(255) NOT NULL,
				PRIMARY KEY (UserID)
			);`

		_, err = db.Exec(createTableSQL)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}

		createTableSQL2 := `
			CREATE TABLE IF NOT EXISTS FollowDB (
				FollowID INT AUTO_INCREMENT,
				OwnerID INT NOT NULL,
				FollowedID INT NOT NULL,
				PRIMARY KEY (FollowID)
				FOREIGN KEY (OwnerID) REFERENCES UserDB(UserID)
				FOREIGN KEY (followedID) REFERENCES UserDB(UserID)
			);`
		_, err = db.Exec(createTableSQL2)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}

		createTableSQL3 := `
			CREATE TABLE IF NOT EXISTS BanDB (
				BanID INT AUTO_INCREMENT,
				OwnerID INT NOT NULL,
				PrayID INT NOT NULL,
				PRIMARY KEY (BanID)
				FOREIGN KEY (OwnerID) REFERENCES UserDB(UserID)
				FOREIGN KEY (PrayID) REFERENCES UserDB(UserID)
			);`
		_, err = db.Exec(createTableSQL3)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}

		createTableSQL6 := `
			CREATE TABLE IF NOT EXISTS PostDB (
				PostID INT AUTO_INCREMENT,
				OwnerID INT NOT NULL,
				PRIMARY KEY (PostID)
				FOREIGN KEY (OwnerID) REFERENCES UserDB(UserID)
			);`
		_, err = db.Exec(createTableSQL6)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}

		createTableSQL4 := `
			CREATE TABLE IF NOT EXISTS LikesDB (
				LikeID INT AUTO_INCREMENT,
				OwnerID INT NOT NULL,
				LikedPhotoID INT NOT NULL,
				PRIMARY KEY (LikeID)
				FOREIGN KEY (OwnerID) REFERENCES UserDB(UserID)
				FOREIGN KEY (LikedPhotoID) REFERENCES PostDB(PostID)
			);`
		_, err = db.Exec(createTableSQL4)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}

		createTableSQL5 := `
			CREATE TABLE IF NOT EXISTS CommentDB (
				CommentID INT AUTO_INCREMENT,
				OwnerID INT NOT NULL,
				PhotoID INT NOT NULL,
				Content VarChar(255) NOT NULL,
				PRIMARY KEY (CommentID)
				FOREIGN KEY (OwnerID) REFERENCES UserDB(UserID)
				FOREIGN KEY (PhotoID) REFERENCES PostDB(PostID)
			);`
		_, err = db.Exec(createTableSQL5)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}

	}

	return &appdbimpl{
		c: db,
	}, nil

}
func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}
