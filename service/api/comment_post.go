package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

type Comment struct {
	Content string `json:"content"`
}

func (rt *_router) CommentPost(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	if r.Method != "POST" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	postIDstr := ps.ByName("postID")
	if postIDstr == "" {
		http.Error(w, "missing postID", http.StatusBadRequest)
		return
	}

	postID, err := strconv.Atoi(postIDstr)
	if err != nil {
		ctx.Logger.Info("Failed to convert postID in int", err.Error())
		http.Error(w, "postID not an int", http.StatusBadRequest)
		return
	}

	ownerIDStr := ps.ByName("ownerID")

	if ownerIDStr == "" {
		http.Error(w, "missing ownerID", http.StatusBadRequest)
		return
	}

	ownerID, err := strconv.Atoi(ownerIDStr)
	if err != nil {
		ctx.Logger.Info("Failed to convert ownerID in int", err.Error())
		http.Error(w, "ownerID not an int", http.StatusBadRequest)
		return
	}

	var comment Comment

	err = json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		ctx.Logger.Info("Failed to decode request body", err.Error())
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if comment.Content == "" {
		http.Error(w, "missing content", http.StatusBadRequest)
		return
	}

	err = rt.db.CommentPost(postID, ownerID, comment.Content)
	if err != nil {
		ctx.Logger.Info("Failed to comment post", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "text/plain")
	_, _ = w.Write([]byte("Post commented!"))
}
