package api

import (
	"encoding/json"
	"net/http"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

type Like struct {
	PostID  int `json:"postID"`
	OwnerID int `json:"ownerID"`
}

func (rt *_router) LikePost(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	var like Like

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

	err = rt.db.LikePost(like.PostID, like.OwnerID)
	if err != nil {
		ctx.Logger.Info("Failed to like post", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "text/plain")
	_, _ = w.Write([]byte("Post Liked!"))
}
