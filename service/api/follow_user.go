package api

import (
	"encoding/json"
	"net/http"

	"github.com/edmii/WasaProject/service/api/reqcontext"
	structs "github.com/edmii/WasaProject/service/models"
	"github.com/julienschmidt/httprouter"
)

// type Follow struct {
// 	FollowedID int `json:"followedID"`
// 	OwnerID    int `json:"ownerID"`
// }

func (rt *_router) FollowUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	var follow structs.Follow

	err := json.NewDecoder(r.Body).Decode(&follow)
	if err != nil {
		ctx.Logger.Info("Failed to decode request body ", err.Error())
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if follow.FollowedID <= 0 {
		http.Error(w, "invalid followedID", http.StatusBadRequest)
		return
	}

	if follow.OwnerID <= 0 {
		http.Error(w, "invalid OwnerID", http.StatusBadRequest)
		return
	}

	result, err := rt.db.FollowUser(follow.OwnerID, follow.FollowedID)
	if err != nil {
		ctx.Logger.Info("Failed to follow user", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Prepare the response
	w.Header().Set("Content-Type", "application/json")

	var response map[string]interface{}

	switch result {
	case 1: // User unfollowed
		response = map[string]interface{}{
			"status":  "success",
			"message": "User unfollowed",
			"data":    follow,
		}
	case 2: // User followed
		response = map[string]interface{}{
			"status":  "success",
			"message": "User followed",
			"data":    follow,
		}
	default:
		http.Error(w, "Unexpected result", http.StatusInternalServerError)
		return
	}

	// Send the response
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		ctx.Logger.Info("Failed to encode response", err.Error())
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (rt *_router) GetFollowed(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var follow structs.Follow

	err := json.NewDecoder(r.Body).Decode(&follow)
	if err != nil {
		ctx.Logger.Info("Failed to decode request body ", err.Error())
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if follow.OwnerID <= 0 {
		http.Error(w, "invalid OwnerID", http.StatusBadRequest)
		return
	}

	followers, err := rt.db.GetFollowed(follow.OwnerID)
	if err != nil {
		ctx.Logger.Info("Failed to get followers", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response := map[string]interface{}{
		"status":   "success",
		"message":  "Followed users retrieved",
		"Followed": followers,
	}

	w.Header().Set("content-type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		ctx.Logger.Info("Failed to encode followers", err.Error())
		http.Error(w, "Failed to encode followers", http.StatusInternalServerError)
		return
	}
}

func (rt *_router) GetFollowers(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var follow structs.Follow

	err := json.NewDecoder(r.Body).Decode(&follow)
	if err != nil {
		ctx.Logger.Info("Failed to decode request body ", err.Error())
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if follow.OwnerID <= 0 {
		http.Error(w, "invalid FollowedID", http.StatusBadRequest)
		return
	}

	followers, err := rt.db.GetFollowers(follow.OwnerID)
	if err != nil {
		ctx.Logger.Info("Failed to get followers", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response := map[string]interface{}{
		"status":    "success",
		"message":   "Following users retrieved",
		"Followers": followers,
	}

	w.Header().Set("content-type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		ctx.Logger.Info("Failed to encode followers", err.Error())
		http.Error(w, "Failed to encode followers", http.StatusInternalServerError)
		return
	}
}
