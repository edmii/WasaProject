package structs

import (
	"time"
)

type FeedResponse struct {
	Username string `json:"username"`
	Posts    []Post `json:"posts"`
}

type User struct {
	Username string `json:"username"`
}

type Follow struct {
	FollowedID int `json:"followedID"`
	OwnerID    int `json:"ownerID"`
}

type Ban struct {
	OwnerID int `json:"ownerID"`
	PrayID  int `json:"prayID"`
}

type Post struct {
	PostID    int    `json:"postID"`
	OwnerID   int    `json:"ownerID"`
	Directory string `json:"imagePath"`
	PostedAt  string `json:"postedAt"`

	RequesterID int `json:"requesterID"`
}

type Like struct {
	PostID  int `json:"postID"`
	OwnerID int `json:"ownerID"`
}

type Comment struct {
	CommentID int       `json:"commentID"`
	Content   string    `json:"content"`
	PostID    int       `json:"postID"`
	OwnerID   int       `json:"ownerID"`
	CreatedAt time.Time `json:"createdAt"`

	RequesterID int `json:"requesterID"`
}
