package api

import (
	"encoding/json"
	"net/http"

	structs "github.com/edmii/WasaProject/service/models"
	"github.com/julienschmidt/httprouter"
)

// type FeedResponse struct {
// 	Username string          `json:"username"`
// 	Posts    []database.Post `json:"posts"`
// }

func (rt *_router) getFeed(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var user structs.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if user.Username == "" {
		http.Error(w, "missing username", http.StatusBadRequest)
		return
	}

	userID, err := rt.db.GetUserID(user.Username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	followers, err := rt.db.GetFollowers(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	followers = append(followers, userID)

	var allPosts []structs.Post

	for _, follower := range followers {
		username := rt.db.GetUsername(follower)
		posts, err := rt.db.GetUserPosts(username)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		allPosts = append(allPosts, posts...)
	}

	response := structs.FeedResponse{
		Username: user.Username,
		Posts:    allPosts,
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

}

func (rt *_router) getProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var user structs.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if user.Username == "" {
		http.Error(w, "missing username", http.StatusBadRequest)
		return
	}

	userID, err := rt.db.GetUserID(user.Username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	posts, err := rt.db.GetUserPosts(user.Username)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	postsCount, err := rt.db.GetPostsCount(userID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	followerCount, err := rt.db.GetFollowersCount(userID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	followedCount, err := rt.db.GetFollowedCount(userID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := structs.ProfilePage{
		Username:      user.Username,
		FollowerCount: followerCount,
		FollowedCount: followedCount,
		Posts:         posts,
		PostsCount:    postsCount,
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return

	}
}
