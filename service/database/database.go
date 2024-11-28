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
	"time"

	structs "github.com/edmii/WasaProject/service/models"
)

// AppDatabase is the high level interface for the DB
type AppDatabase interface {
	GetName() (string, error)
	SetName(name string) error
	GetDatabaseTableContent(tableName string) ([]map[string]interface{}, error)
	DestroyDB() error

	CreateUser(username string) error
	GetUserID(username string) (int, error)
	GetUsername(userID int) string

	CreatePost(ownerID int, directory string, PostedAt time.Time) (int, error)
	DeletePost(postID int, requesterID int) error
	GetUserPosts(username string) ([]structs.Post, error)
	GetPostsCount(userID int) (int, error)

	LikePost(PostID int, OwnerID int) (int, error)
	GetLikes(ownerID int) ([]int, error)
	GetLikesCount(postID int) (int, error)

	CommentPost(PostID int, OwnerID int, Content string, CreatedAt time.Time) error
	DeleteComment(CommentID int, OwnerID int, PostID int) error
	GetComments(postID int) ([]structs.Comment, error)
	GetCommentsCount(postID int) (int, error)

	BanUser(OwnerID int, PrayID int) (int, error)
	GetBannedUsers(ownerID int) ([]int, error)

	FollowUser(OwnerID int, FollowedID int) (int, error)
	GetFollowed(ownerID int) ([]int, error)
	GetFollowers(ownerID int) ([]int, error)

	// FollowUser(ownerID int, followedID int) error
	// UnfollowUser(ownerID int, followedID int) error
	// BanUser(ownerID int, prayID int) error
	// UnbanUser(ownerID int, prayID int) error

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

	// Enable foreign key constraints for SQLite
	_, err := db.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		return nil, fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	// Create the User table
	createTableSQL := `
			CREATE TABLE IF NOT EXISTS UserDB (
				UserID INTEGER NOT NULL PRIMARY KEY,
				Username VARCHAR(255) NOT NULL UNIQUE
			);`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		return nil, fmt.Errorf("error creating database structure: %w", err)
	}

	createTableSQL2 := `
			CREATE TABLE IF NOT EXISTS FollowDB (
				FollowID INTEGER NOT NULL PRIMARY KEY,
				OwnerID INT NOT NULL,
				FollowedID INT NOT NULL,
				FOREIGN KEY (OwnerID) REFERENCES UserDB(UserID),
				FOREIGN KEY (followedID) REFERENCES UserDB(UserID)
			);`
	_, err = db.Exec(createTableSQL2)
	if err != nil {
		return nil, fmt.Errorf("error creating database structure: %w", err)
	}

	createTableSQL3 := `
			CREATE TABLE IF NOT EXISTS BanDB (
				BanID INTEGER NOT NULL PRIMARY KEY,
				OwnerID INT NOT NULL,
				PrayID INT NOT NULL,
				FOREIGN KEY (OwnerID) REFERENCES UserDB(UserID),
				FOREIGN KEY (PrayID) REFERENCES UserDB(UserID)
			);`
	_, err = db.Exec(createTableSQL3)
	if err != nil {
		return nil, fmt.Errorf("error creating database structure: %w", err)
	}

	createTableSQL6 := `
			CREATE TABLE IF NOT EXISTS PostDB (
				PostID INTEGER NOT NULL PRIMARY KEY,
				Directory STRING NOT NULL,
				OwnerID INT NOT NULL,
				PostedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				FOREIGN KEY (OwnerID) REFERENCES UserDB(UserID)
			);`
	_, err = db.Exec(createTableSQL6)
	if err != nil {
		return nil, fmt.Errorf("error creating database structure: %w", err)
	}

	createTableSQL4 := `
			CREATE TABLE IF NOT EXISTS LikesDB (
				LikeID INTEGER NOT NULL PRIMARY KEY,
				OwnerID INT NOT NULL,
				LikedPhotoID INT NOT NULL,
				FOREIGN KEY (OwnerID) REFERENCES UserDB(UserID),
				FOREIGN KEY (LikedPhotoID) REFERENCES PostDB(PostID)
			);`
	_, err = db.Exec(createTableSQL4)
	if err != nil {
		return nil, fmt.Errorf("error creating database structure: %w", err)
	}

	createTableSQL5 := `
			CREATE TABLE IF NOT EXISTS CommentDB (
				CommentID INTEGER NOT NULL PRIMARY KEY,
				OwnerID INT NOT NULL,
				PhotoID INT NOT NULL,
				Content VarChar(255) NOT NULL,
				CreatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				FOREIGN KEY (OwnerID) REFERENCES UserDB(UserID),
				FOREIGN KEY (PhotoID) REFERENCES PostDB(PostID)
			);`
	_, err = db.Exec(createTableSQL5)
	if err != nil {
		return nil, fmt.Errorf("error creating database structure: %w", err)
	}

	return &appdbimpl{
		c: db,
	}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}
