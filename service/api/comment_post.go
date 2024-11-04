package api

import (
	"encoding/json"
	"net/http"
	"time"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

type Comment struct {
	Content string    `json:"content"`
	PostID  int       `json:"postID"`
	OwnerID int       `json:"ownerID"`
	Created time.Time `json:"created"`
}

func (rt *_router) CommentPost(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	if r.Method != "POST" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var comment Comment

	comment.Created = time.Now()

	err := json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		ctx.Logger.Info("Failed to decode request body ", err.Error())
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if comment.Content == "" {
		http.Error(w, "missing content", http.StatusBadRequest)
		return
	}

	if comment.PostID <= 0 {
		http.Error(w, "invalid post ID", http.StatusBadRequest)
		return
	}

	if comment.OwnerID <= 0 {
		http.Error(w, "invalid ownerID", http.StatusBadRequest)
		return
	}

	err = rt.db.CommentPost(comment.PostID, comment.OwnerID, comment.Content, comment.Created)
	if err != nil {
		ctx.Logger.Info("Failed to comment post", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "text/plain")
	_, _ = w.Write([]byte("Post commented!"))

}

func (rt *_router) GetComments(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var comment Comment

	err := json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		ctx.Logger.Info("Failed to decode request body ", err.Error())
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if comment.PostID <= 0 {
		http.Error(w, "invalid post ID", http.StatusBadRequest)
		return
	}

	comments, err := rt.db.GetComments(comment.PostID)
	if err != nil {
		ctx.Logger.Info("Failed to get comments", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "application/json")
	_ = json.NewEncoder(w).Encode(comments)
}
