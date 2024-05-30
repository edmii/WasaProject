package api

import (
	"encoding/json"
	"net/http"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
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

	err = rt.db.FollowUser(follow.OwnerID, follow.FollowedID)
	if err != nil {
		ctx.Logger.Info("Failed to follow user", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "text/plain")
	_, _ = w.Write([]byte("User followed!"))
}
