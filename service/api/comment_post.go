package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/edmii/WasaProject/service/api/reqcontext"
	structs "github.com/edmii/WasaProject/service/models"
	"github.com/edmii/WasaProject/service/utils"
	"github.com/julienschmidt/httprouter"
)

// type Comment struct {
// 	CommentID int       `json:"commentID"`
// 	Content   string    `json:"content"`
// 	PostID    int       `json:"postID"`
// 	OwnerID   int       `json:"ownerID"`
// 	CreatedAt time.Time `json:"createdAt"`

// 	RequesterID int `json:"requesterID"`
// }

func (rt *_router) CommentPost(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	if r.Method != "POST" {
		utils.SendErrorResponse(w, "Invalid request method", []string{"Only POST requests are allowed"}, http.StatusMethodNotAllowed)
		return
	}

	var comment structs.Comment
	comment.CreatedAt = time.Now()

	err := json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		ctx.Logger.Info("Failed to decode request body ", err.Error())
		utils.SendErrorResponse(w, "Invalid request body", []string{"Failed to decode JSON request body"}, http.StatusBadRequest)
		return
	}

	if comment.Content == "" {
		utils.SendErrorResponse(w, "Invalid request", []string{"Missing content"}, http.StatusBadRequest)
		return
	}

	if comment.PostID <= 0 {
		utils.SendErrorResponse(w, "Invalid request", []string{"Invalid PostID"}, http.StatusBadRequest)
		return
	}

	if comment.OwnerID <= 0 {
		utils.SendErrorResponse(w, "Invalid request", []string{"Invalid OwnerID"}, http.StatusBadRequest)
		return
	}

	err = rt.db.CommentPost(comment.PostID, comment.OwnerID, comment.Content, comment.CreatedAt)
	if err != nil {
		ctx.Logger.Info("Failed to comment post", err.Error())
		utils.SendErrorResponse(w, "Database error", []string{"Failed to post comment", err.Error()}, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	response := map[string]interface{}{
		"status":  "success",
		"message": "Comment posted",
		"data":    comment,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		ctx.Logger.Info("Failed to encode response", err.Error())
		utils.SendErrorResponse(w, "Server error", []string{"Failed to encode response", err.Error()}, http.StatusInternalServerError)
	}
}

// func (rt *_router) CommentPost(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
// 	if r.Method != "POST" {
// 		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
// 		return
// 	}

// 	var comment structs.Comment

// 	comment.CreatedAt = time.Now()

// 	err := json.NewDecoder(r.Body).Decode(&comment)
// 	if err != nil {
// 		ctx.Logger.Info("Failed to decode request body ", err.Error())
// 		http.Error(w, "Invalid request body", http.StatusBadRequest)
// 		return
// 	}

// 	if comment.Content == "" {
// 		http.Error(w, "missing content", http.StatusBadRequest)
// 		return
// 	}

// 	if comment.PostID <= 0 {
// 		http.Error(w, "invalid post ID", http.StatusBadRequest)
// 		return
// 	}

// 	if comment.OwnerID <= 0 {
// 		http.Error(w, "invalid ownerID", http.StatusBadRequest)
// 		return
// 	}

// 	err = rt.db.CommentPost(comment.PostID, comment.OwnerID, comment.Content, comment.CreatedAt)
// 	if err != nil {
// 		ctx.Logger.Info("Failed to comment post", err.Error())
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("content-type", "application/json")

// 	response := map[string]interface{}{
// 		"status":  "success",
// 		"message": "Comment posted",
// 		"data":    comment,
// 	}

//		err = json.NewEncoder(w).Encode(response)
//		if err != nil {
//			ctx.Logger.Info("Failed to encode posts", err.Error())
//			http.Error(w, "Failed to encode posts", http.StatusInternalServerError)
//			return
//		}
//	}
func (rt *_router) DeleteComment(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	if r.Method != "DELETE" {
		utils.SendErrorResponse(w, "Invalid request method", []string{"Only DELETE requests are allowed"}, http.StatusMethodNotAllowed)
		return
	}

	var comment structs.Comment

	err := json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		ctx.Logger.Info("Failed to decode request body ", err.Error())
		utils.SendErrorResponse(w, "Invalid request body", []string{"Failed to decode JSON request body"}, http.StatusBadRequest)
		return
	}

	if comment.CommentID <= 0 {
		utils.SendErrorResponse(w, "Invalid request", []string{"Invalid CommentID"}, http.StatusBadRequest)
		return
	}

	err = rt.db.DeleteComment(comment.CommentID, comment.RequesterID, comment.PostID)
	if err != nil {
		ctx.Logger.Info("Failed to delete comment", err.Error())
		utils.SendErrorResponse(w, "Database error", []string{"Failed to delete comment", err.Error()}, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	_, _ = w.Write([]byte("Comment deleted"))
}

// func (rt *_router) DeleteComment(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
// 	if r.Method != "DELETE" {
// 		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
// 		return
// 	}

// 	var comment structs.Comment

// 	err := json.NewDecoder(r.Body).Decode(&comment)
// 	if err != nil {
// 		ctx.Logger.Info("Failed to decode request body ", err.Error())
// 		http.Error(w, "Invalid request body", http.StatusBadRequest)
// 		return
// 	}

// 	if comment.CommentID <= 0 {
// 		http.Error(w, "invalid post ID", http.StatusBadRequest)
// 		return
// 	}

// 	err = rt.db.DeleteComment(comment.CommentID, comment.RequesterID, comment.PostID)
// 	if err != nil {
// 		ctx.Logger.Info("Failed to delete comment", err.Error())
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("content-type", "text/plain")
// 	_, _ = w.Write([]byte("Comment deleted"))
// }

func (rt *_router) GetComments(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var comment structs.Comment

	err := json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		ctx.Logger.Info("Failed to decode request body ", err.Error())
		utils.SendErrorResponse(w, "Invalid request body", []string{"Failed to decode JSON request body"}, http.StatusBadRequest)
		return
	}

	if comment.PostID <= 0 {
		utils.SendErrorResponse(w, "Invalid request", []string{"Invalid PostID"}, http.StatusBadRequest)
		return
	}

	comments, err := rt.db.GetComments(comment.PostID)
	if err != nil {
		ctx.Logger.Info("Failed to get comments", err.Error())
		utils.SendErrorResponse(w, "Database error", []string{"Failed to retrieve comments for the post", err.Error()}, http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"status":  "success",
		"message": "Comments retrieved successfully",
		"data":    comments,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		ctx.Logger.Info("Failed to encode comments", err.Error())
		utils.SendErrorResponse(w, "Server error", []string{"Failed to encode response", err.Error()}, http.StatusInternalServerError)
	}
}

// func (rt *_router) GetComments(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
// 	var comment structs.Comment

// 	err := json.NewDecoder(r.Body).Decode(&comment)
// 	if err != nil {
// 		ctx.Logger.Info("Failed to decode request body ", err.Error())
// 		http.Error(w, "Invalid request body", http.StatusBadRequest)
// 		return
// 	}

// 	if comment.PostID <= 0 {
// 		http.Error(w, "invalid post ID", http.StatusBadRequest)
// 		return
// 	}

// 	comments, err := rt.db.GetComments(comment.PostID)
// 	if err != nil {
// 		ctx.Logger.Info("Failed to get comments", err.Error())
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("content-type", "application/json")

// 	response := map[string]interface{}{
// 		"status":  "success",
// 		"message": "User posts retrieved",
// 		"data":    comments,
// 	}

// 	err = json.NewEncoder(w).Encode(response)
// 	if err != nil {
// 		ctx.Logger.Info("Failed to encode posts", err.Error())
// 		http.Error(w, "Failed to encode posts", http.StatusInternalServerError)
// 		return
// 	}
// }
