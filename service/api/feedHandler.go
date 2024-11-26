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

	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

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
