package api

import (
	"encoding/json"
	"net/http"

	structs "github.com/edmii/WasaProject/service/models"
	"github.com/edmii/WasaProject/service/utils"
	"github.com/julienschmidt/httprouter"
)

// type FeedResponse struct {
// 	Username string          `json:"username"`
// 	Posts    []database.Post `json:"posts"`
// }

// func (rt *_router) getFeed(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
// 	var user structs.User
// 	err := json.NewDecoder(r.Body).Decode(&user)
// 	if err != nil {
// 		http.Error(w, "Invalid request body", http.StatusBadRequest)
// 		return
// 	}

// 	if user.Username == "" {
// 		http.Error(w, "missing username", http.StatusBadRequest)
// 		return
// 	}

// 	userID, err := rt.db.GetUserID(user.Username)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	followers, err := rt.db.GetFollowers(userID)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	followers = append(followers, userID)

// 	var allPosts []structs.Post

// 	for _, follower := range followers {
// 		username := rt.db.GetUsername(follower)
// 		posts, err := rt.db.GetUserPosts(username)
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}

// 		allPosts = append(allPosts, posts...)
// 	}

// 	feed := structs.FeedResponse{
// 		Username: user.Username,
// 		Posts:    allPosts,
// 	}

// 	w.Header().Set("Content-Type", "application/json")

// 	response := map[string]interface{}{
// 		"status":  "success",
// 		"message": "Feed retrieved",
// 		"data":    feed,
// 	}

// 	err = json.NewEncoder(w).Encode(response)
// 	if err != nil {
// 		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
// 		return
// 	}
// }

func (rt *_router) getFeed(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var user structs.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		utils.SendErrorResponse(w, "Invalid request body", []string{"Failed to decode JSON request body"}, http.StatusBadRequest)
		return
	}

	if user.Username == "" {
		utils.SendErrorResponse(w, "Invalid request body", []string{"Missing username in request body"}, http.StatusBadRequest)
		return
	}

	userID, err := rt.db.GetUserID(user.Username)
	if err != nil {
		utils.SendErrorResponse(w, "Database error", []string{"Failed to retrieve user ID", err.Error()}, http.StatusInternalServerError)
		return
	}

	followers, err := rt.db.GetFollowers(userID)
	if err != nil {
		utils.SendErrorResponse(w, "Database error", []string{"Failed to retrieve followers", err.Error()}, http.StatusInternalServerError)
		return
	}

	followers = append(followers, userID)

	var allPosts []structs.Post
	for _, follower := range followers {
		username := rt.db.GetUsername(follower)
		posts, err := rt.db.GetUserPosts(username)
		if err != nil {
			utils.SendErrorResponse(w, "Database error", []string{"Failed to retrieve posts for user", username, err.Error()}, http.StatusInternalServerError)
			return
		}
		allPosts = append(allPosts, posts...)
	}

	feed := structs.FeedResponse{
		Username: user.Username,
		Posts:    allPosts,
	}

	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"status":  "success",
		"message": "Feed retrieved",
		"data":    feed,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		utils.SendErrorResponse(w, "Server error", []string{"Failed to encode response", err.Error()}, http.StatusInternalServerError)
	}
}

func (rt *_router) getProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var user structs.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		utils.SendErrorResponse(w, "Invalid request body", []string{"Failed to decode JSON request body"}, http.StatusBadRequest)
		return
	}

	// Check username
	if user.Username == "" {
		utils.SendErrorResponse(w, "Invalid request body", []string{"Missing username in request body"}, http.StatusBadRequest)
		return
	}

	// Check requester ID
	if user.ID <= 0 {
		utils.SendErrorResponse(w, "Invalid request body", []string{"Invalid requester ID"}, http.StatusBadRequest)
		return
	}

	userID, err := rt.db.GetUserID(user.Username)
	if err != nil {
		utils.SendErrorResponse(w, "Database error", []string{"Failed to retrieve user ID", err.Error()}, http.StatusInternalServerError)
		return
	}

	// Check ban status
	banExists, err := rt.db.CheckBanStatus(userID, user.ID)
	if err != nil {
		utils.SendErrorResponse(w, "Database error", []string{"Failed to check ban status", err.Error()}, http.StatusInternalServerError)
		return
	}
	if banExists {
		utils.SendErrorResponse(w, "Access denied", []string{"You are banned from viewing this profile"}, http.StatusForbidden)
		return
	}

	posts, err := rt.db.GetUserPosts(user.Username)
	if err != nil {
		utils.SendErrorResponse(w, "Database error", []string{"Failed to retrieve user posts", err.Error()}, http.StatusInternalServerError)
		return
	}

	postsCount, err := rt.db.GetPostsCount(userID)
	if err != nil {
		utils.SendErrorResponse(w, "Database error", []string{"Failed to retrieve post count", err.Error()}, http.StatusInternalServerError)
		return
	}

	followerCount, err := rt.db.GetFollowersCount(userID)
	if err != nil {
		utils.SendErrorResponse(w, "Database error", []string{"Failed to retrieve follower count", err.Error()}, http.StatusInternalServerError)
		return
	}

	followedCount, err := rt.db.GetFollowedCount(userID)
	if err != nil {
		utils.SendErrorResponse(w, "Database error", []string{"Failed to retrieve followed count", err.Error()}, http.StatusInternalServerError)
		return
	}

	profile := structs.ProfilePage{
		Username:      user.Username,
		FollowerCount: followerCount,
		FollowedCount: followedCount,
		Posts:         posts,
		PostsCount:    postsCount,
	}

	w.Header().Set("Content-Type", "application/json")

	response := map[string]interface{}{
		"status":  "success",
		"message": "Profile page retrieved",
		"data":    profile,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		utils.SendErrorResponse(w, "Server error", []string{"Failed to encode response", err.Error()}, http.StatusInternalServerError)
	}
}

// func (rt *_router) getProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
// 	var user structs.User
// 	err := json.NewDecoder(r.Body).Decode(&user)
// 	if err != nil {
// 		http.Error(w, "Invalid request body", http.StatusBadRequest)
// 		return
// 	}

// 	//username of profile getting retrieved
// 	if user.Username == "" {
// 		http.Error(w, "missing username", http.StatusBadRequest)
// 		return
// 	}

// 	//id of requester
// 	if user.ID <= 0 {
// 		http.Error(w, "invalid requester ID", http.StatusBadRequest)
// 		return
// 	}

// 	userID, err := rt.db.GetUserID(user.Username)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	//check ban existance
// 	banExists, err := rt.db.CheckBanStatus(userID, user.ID)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	if banExists {
// 		return
// 	}

// 	posts, err := rt.db.GetUserPosts(user.Username)

// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	postsCount, err := rt.db.GetPostsCount(userID)

// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	followerCount, err := rt.db.GetFollowersCount(userID)

// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	followedCount, err := rt.db.GetFollowedCount(userID)

// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	profile := structs.ProfilePage{
// 		Username:      user.Username,
// 		FollowerCount: followerCount,
// 		FollowedCount: followedCount,
// 		Posts:         posts,
// 		PostsCount:    postsCount,
// 	}

// 	w.Header().Set("Content-Type", "application/json")

// 	response := map[string]interface{}{
// 		"status":  "success",
// 		"message": "Profile page retrieved",
// 		"data":    profile,
// 	}

// 	err = json.NewEncoder(w).Encode(response)
// 	if err != nil {
// 		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
// 		return

// 	}
// }

// func sendErrorResponse(w http.ResponseWriter, message string, details []string, statusCode int) {
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(statusCode)

// 	response := structs.ErrorResponse{
// 		Status:  fmt.Sprintf("error %d", statusCode),
// 		Message: message,
// 		Details: details,
// 	}

// 	_ = json.NewEncoder(w).Encode(response)
// }
