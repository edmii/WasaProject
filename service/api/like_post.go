package api

import (
	"encoding/json"
	"net/http"

	"github.com/edmii/WasaProject/service/api/reqcontext"
	structs "github.com/edmii/WasaProject/service/models"
	"github.com/julienschmidt/httprouter"
)

// type Like struct {
// 	PostID  int `json:"postID"`
// 	OwnerID int `json:"ownerID"`
// }

func (rt *_router) LikePost(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	var like structs.Like

	err := json.NewDecoder(r.Body).Decode(&like)
	if err != nil {
		ctx.Logger.Info("Failed to decode request body ", err.Error())
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if like.PostID <= 0 {
		http.Error(w, "invalid post ID", http.StatusBadRequest)
		return
	}

	if like.OwnerID <= 0 {
		http.Error(w, "invalid ownerID", http.StatusBadRequest)
		return
	}

	result, err := rt.db.LikePost(like.PostID, like.OwnerID)
	if err != nil {
		ctx.Logger.Info("Failed to like post", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Prepare the response
	w.Header().Set("Content-Type", "application/json")

	var response map[string]interface{}

	switch result {
	case 1: // Post unliked
		response = map[string]interface{}{
			"status":  "success",
			"message": "Post unliked",
			"data":    like,
		}
	case 2: // Post liked
		response = map[string]interface{}{
			"status":  "success",
			"message": "Post liked",
			"data":    like,
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

func (rt *_router) GetLikes(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var like structs.Like

	err := json.NewDecoder(r.Body).Decode(&like)
	if err != nil {
		ctx.Logger.Info("Failed to decode request body ", err.Error())
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if like.PostID <= 0 {
		http.Error(w, "invalid post ID", http.StatusBadRequest)
		return
	}

	likes, err := rt.db.GetLikes(like.PostID)
	if err != nil {
		ctx.Logger.Info("Failed to get likes", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "application/json")
	err = json.NewEncoder(w).Encode(likes)
	if err != nil {
		ctx.Logger.Info("Failed to encode likes", err.Error())
		http.Error(w, "Failed to encode likes", http.StatusInternalServerError)
		return
	}
}
