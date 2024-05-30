package api

import (
	"encoding/json"
	"net/http"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

type Unban struct {
	OwnerID int `json:"ownerID"`
	PrayID  int `json:"prayID"`
}

func (rt *_router) UnbanUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	var unban Unban

	err := json.NewDecoder(r.Body).Decode(&unban)
	if err != nil {
		ctx.Logger.Info("Failed to decode request body ", err.Error())
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if unban.OwnerID <= 0 {
		http.Error(w, "Invalid OwnerID", http.StatusBadRequest)
		return
	}

	if unban.PrayID <= 0 {
		http.Error(w, "Invalid PrayID", http.StatusBadRequest)
		return
	}

	err = rt.db.UnbanUser(unban.OwnerID, unban.PrayID)
	if err != nil {
		ctx.Logger.Info("Failed to unban user", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "text/plain")
	_, _ = w.Write([]byte("User unbanned!"))
}
