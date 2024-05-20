package api

import (
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) LikePost(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	postIDstr := ps.ByName("postID")
	if postIDstr == "" {
		http.Error(w, "missing post ID", http.StatusBadRequest)
		return
	}

	postID, err := strconv.Atoi(postIDstr)
	if err != nil {
		http.Error(w, "postID not an int", http.StatusBadRequest)
		return
	}

	ownerIDstr := ps.ByName("ownerID")

	if ownerIDstr == "" {
		http.Error(w, "missing ownerID", http.StatusBadRequest)
		return
	}

	ownerID, err := strconv.Atoi(ownerIDstr)
	if err != nil {
		http.Error(w, "ownerID not an int", http.StatusBadRequest)
		return
	}

	err = rt.db.LikePost(postID, ownerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "text/plain")
	_, _ = w.Write([]byte("Post Liked!"))
}
