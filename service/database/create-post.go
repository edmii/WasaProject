package database

type Post struct {
	PostID    int    `json:"postID"`
	OwnerID   int    `json:"ownerID"`
	Directory string `json:"imagePath"`
	PostedAt  string `json:"postedAt"`
}

func (db *appdbimpl) CreatePost(ownerID int, directory string) error {

	insertQuery := "INSERT INTO PostDB (OwnerID, Directory) VALUES ($1, $2)"

	_, err := db.c.Exec(insertQuery, ownerID, directory)

	if err != nil {
		return err
	}

	return nil
}

func (db *appdbimpl) GetUserPosts(username string) ([]Post, error) {

	//("SELECT UserID FROM UserDB WHERE Username = $1", username)
	rows, err := db.c.Query("SELECT * FROM PostDB WHERE OwnerID = (SELECT UserID FROM UserDB WHERE Username = $1)", username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	//return a slice of json, with all info about the post
	var posts []Post
	for rows.Next() {
		var post Post
		if err := rows.Scan(&post.PostID, &post.OwnerID, &post.Directory, &post.PostedAt); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil

}
