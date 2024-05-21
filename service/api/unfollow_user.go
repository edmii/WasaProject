package api

import (
	"net/http"
	"strconv"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) UnfollowUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	ownerIdstr := ps.ByName("ownerID")

	if ownerIdstr == "" {
		http.Error(w, "missing ownerID", http.StatusBadRequest)
		return
	}

	ownerID, err := strconv.Atoi(ownerIdstr)
	if err != nil {
		ctx.Logger.Info("Failed to convert OwnerID in int", err.Error())
		http.Error(w, "ownerID not an int", http.StatusBadRequest)
		return
	}

	followedIDStr := ps.ByName("followedID")

	if followedIDStr == "" {
		ctx.Logger.Info("Failed to get followedID", err.Error())
		http.Error(w, "missing followedID", http.StatusBadRequest)
		return
	}

	followedID, err := strconv.Atoi(followedIDStr)
	if err != nil {
		ctx.Logger.Info("Failed to convert followedID in int", err.Error())
		http.Error(w, "followedID not an int", http.StatusBadRequest)
		return
	}

	err = rt.db.UnfollowUser(ownerID, followedID)
	if err != nil {
		ctx.Logger.Info("Failed to unfollow user", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "text/plain")
	_, _ = w.Write([]byte("User unfollowed!"))
}
