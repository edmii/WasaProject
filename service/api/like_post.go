package api

import (
	"encoding/json"
	"net/http"

	"github.com/edmii/WasaProject/service/api/reqcontext"
	structs "github.com/edmii/WasaProject/service/models"
	"github.com/edmii/WasaProject/service/utils"
	"github.com/julienschmidt/httprouter"
)

// type Like struct {
// 	PostID  int `json:"postID"`
// 	OwnerID int `json:"ownerID"`
// }

func (rt *_router) likePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var like structs.Like

	err := json.NewDecoder(r.Body).Decode(&like)
	if err != nil {
		ctx.Logger.Info("Failed to decode request body ", err.Error())
		utils.SendErrorResponse(w, "Invalid request body", []string{"Failed to decode JSON request body"}, http.StatusBadRequest)
		return
	}

	if like.PostID <= 0 {
		utils.SendErrorResponse(w, "Invalid request", []string{"Invalid PostID"}, http.StatusBadRequest)
		return
	}

	if like.OwnerID <= 0 {
		utils.SendErrorResponse(w, "Invalid request", []string{"Invalid OwnerID"}, http.StatusBadRequest)
		return
	}

	result, err := rt.db.LikePost(like.PostID, like.OwnerID)
	if err != nil {
		ctx.Logger.Info("Failed to like post", err.Error())
		utils.SendErrorResponse(w, "Database error", []string{"Failed to like/unlike post", err.Error()}, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	var response map[string]interface{}

	switch result {
	case 1: // do nothing
		response = map[string]interface{}{
			"status":  "success",
			"message": "Post already liked",
			"data":    like,
		}
	case 2: // Post liked
		response = map[string]interface{}{
			"status":  "success",
			"message": "Post liked",
			"data":    like,
		}
	default:
		utils.SendErrorResponse(w, "Unexpected result", []string{"Unknown result code from LikePost operation"}, http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		ctx.Logger.Info("Failed to encode response", err.Error())
		utils.SendErrorResponse(w, "Server error", []string{"Failed to encode response", err.Error()}, http.StatusInternalServerError)
	}
}

func (rt *_router) unlikePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var like structs.Like

	err := json.NewDecoder(r.Body).Decode(&like)
	if err != nil {
		ctx.Logger.Info("Failed to decode request body ", err.Error())
		utils.SendErrorResponse(w, "Invalid request body", []string{"Failed to decode JSON request body"}, http.StatusBadRequest)
		return
	}

	if like.PostID <= 0 {
		utils.SendErrorResponse(w, "Invalid request", []string{"Invalid PostID"}, http.StatusBadRequest)
		return
	}

	if like.OwnerID <= 0 {
		utils.SendErrorResponse(w, "Invalid request", []string{"Invalid OwnerID"}, http.StatusBadRequest)
		return
	}

	result, err := rt.db.UnlikePost(like.PostID, like.OwnerID)
	if err != nil {
		ctx.Logger.Info("Failed to like post", err.Error())
		utils.SendErrorResponse(w, "Database error", []string{"Failed to like/unlike post", err.Error()}, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	var response map[string]interface{}

	switch result {
	case 1: // Post unliked
		response = map[string]interface{}{
			"status":  "success",
			"message": "Post unliked",
			"data":    like,
		}
	case 2: // do nothing
		response = map[string]interface{}{
			"status":  "success",
			"message": "Post already not liked",
			"data":    like,
		}
	default:
		utils.SendErrorResponse(w, "Unexpected result", []string{"Unknown result code from LikePost operation"}, http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		ctx.Logger.Info("Failed to encode response", err.Error())
		utils.SendErrorResponse(w, "Server error", []string{"Failed to encode response", err.Error()}, http.StatusInternalServerError)
	}
}

// func (rt *_router) LikePost(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

// 	var like structs.Like

// 	err := json.NewDecoder(r.Body).Decode(&like)
// 	if err != nil {
// 		ctx.Logger.Info("Failed to decode request body ", err.Error())
// 		http.Error(w, "Invalid request body", http.StatusBadRequest)
// 		return
// 	}

// 	if like.PostID <= 0 {
// 		http.Error(w, "invalid post ID", http.StatusBadRequest)
// 		return
// 	}

// 	if like.OwnerID <= 0 {
// 		http.Error(w, "invalid ownerID", http.StatusBadRequest)
// 		return
// 	}

// 	result, err := rt.db.LikePost(like.PostID, like.OwnerID)
// 	if err != nil {
// 		ctx.Logger.Info("Failed to like post", err.Error())
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	// Prepare the response
// 	w.Header().Set("Content-Type", "application/json")

// 	var response map[string]interface{}

// 	switch result {
// 	case 1: // Post unliked
// 		response = map[string]interface{}{
// 			"status":  "success",
// 			"message": "Post unliked",
// 			"data":    like,
// 		}
// 	case 2: // Post liked
// 		response = map[string]interface{}{
// 			"status":  "success",
// 			"message": "Post liked",
// 			"data":    like,
// 		}
// 	default:
// 		http.Error(w, "Unexpected result", http.StatusInternalServerError)
// 		return
// 	}

//		// Send the response
//		err = json.NewEncoder(w).Encode(response)
//		if err != nil {
//			ctx.Logger.Info("Failed to encode response", err.Error())
//			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
//		}
//	}
func (rt *_router) getLikes(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var like structs.Like

	err := json.NewDecoder(r.Body).Decode(&like)
	if err != nil {
		ctx.Logger.Info("Failed to decode request body ", err.Error())
		utils.SendErrorResponse(w, "Invalid request body", []string{"Failed to decode JSON request body"}, http.StatusBadRequest)
		return
	}

	if like.PostID <= 0 {
		utils.SendErrorResponse(w, "Invalid request", []string{"Invalid PostID"}, http.StatusBadRequest)
		return
	}

	likes, err := rt.db.GetLikes(like.PostID)
	if err != nil {
		ctx.Logger.Info("Failed to get likes", err.Error())
		utils.SendErrorResponse(w, "Database error", []string{"Failed to retrieve likes for the post", err.Error()}, http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"status":  "success",
		"message": "Likes retrieved successfully",
		"data":    likes,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		ctx.Logger.Info("Failed to encode likes", err.Error())
		utils.SendErrorResponse(w, "Server error", []string{"Failed to encode response", err.Error()}, http.StatusInternalServerError)
	}
}

// func (rt *_router) GetLikes(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
// 	var like structs.Like

// 	err := json.NewDecoder(r.Body).Decode(&like)
// 	if err != nil {
// 		ctx.Logger.Info("Failed to decode request body ", err.Error())
// 		http.Error(w, "Invalid request body", http.StatusBadRequest)
// 		return
// 	}

// 	if like.PostID <= 0 {
// 		http.Error(w, "invalid post ID", http.StatusBadRequest)
// 		return
// 	}

// 	likes, err := rt.db.GetLikes(like.PostID)
// 	if err != nil {
// 		ctx.Logger.Info("Failed to get likes", err.Error())
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("content-type", "application/json")
// 	err = json.NewEncoder(w).Encode(likes)
// 	if err != nil {
// 		ctx.Logger.Info("Failed to encode likes", err.Error())
// 		http.Error(w, "Failed to encode likes", http.StatusInternalServerError)
// 		return
// 	}
// }
