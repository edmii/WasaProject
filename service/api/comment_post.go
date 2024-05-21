package api

import (
	"net/http"
	"strconv"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) CommentPost(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	if r.Method != "POST" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

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

	ownerIDStr := ps.ByName("ownerID")

	if ownerIDStr == "" {
		http.Error(w, "missing ownerID", http.StatusBadRequest)
		return
	}

	ownerID, err := strconv.Atoi(ownerIDStr)
	if err != nil {
		http.Error(w, "ownerID not an int", http.StatusBadRequest)
		return
	}

	content := ps.ByName("content")

	if content == "" {
		http.Error(w, "missing ownerID", http.StatusBadRequest)
		return
	}

	err = rt.db.CommentPost(postID, ownerID, content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "text/plain")
	_, _ = w.Write([]byte("Post commented!"))
}
