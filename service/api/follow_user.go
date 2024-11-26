package api

import (
	"encoding/json"
	"net/http"

	"github.com/edmii/WasaProject/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

type Follow struct {
	FollowedID int `json:"followedID"`
	OwnerID    int `json:"ownerID"`
}

func (rt *_router) FollowUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	var follow Follow

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

	w.Header().Set("content-type", "text/plain")
	switch result {
	case 1:
		_, _ = w.Write([]byte("User unfollowed!"))
	case 2:
		_, _ = w.Write([]byte("User followed!"))
	}
}

func (rt *_router) GetFollowed(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var follow Follow

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
	response := map[string][]int{
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
	var follow Follow

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
	response := map[string][]int{
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
