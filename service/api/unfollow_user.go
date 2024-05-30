package api

import (
	"encoding/json"
	"net/http"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

type Unfollow struct {
	UnfollowedID int `json:"unfollowedID"`
	OwnerID      int `json:"ownerID"`
}

func (rt *_router) UnfollowUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	var unfollow Unfollow

	err := json.NewDecoder(r.Body).Decode(&unfollow)
	if err != nil {
		ctx.Logger.Info("Failed to decode request body ", err.Error())
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if unfollow.OwnerID <= 0 {
		http.Error(w, "ownerID not valid", http.StatusBadRequest)
		return
	}

	if unfollow.UnfollowedID <= 0 {
		http.Error(w, "unfollowedID not valid", http.StatusBadRequest)
		return
	}

	err = rt.db.UnfollowUser(unfollow.OwnerID, unfollow.UnfollowedID)
	if err != nil {
		ctx.Logger.Info("Failed to unfollow user", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "text/plain")
	_, _ = w.Write([]byte("User unfollowed!"))
}
