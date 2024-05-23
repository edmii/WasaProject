package api

import (
	"encoding/json"
	"net/http"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

type Ban struct {
	OwnerID int `json:"ownerID"`
	PrayID  int `json:"prayID"`
}

func (rt *_router) BanUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	var ban Ban

	err := json.NewDecoder(r.Body).Decode(&ban)
	if err != nil {
		ctx.Logger.Info("Failed to decode request body ", err.Error())
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if ban.OwnerID <= 0 {
		http.Error(w, "Invalid OwnerID", http.StatusBadRequest)
		return
	}

	if ban.PrayID <= 0 {
		http.Error(w, "Invalid PrayID", http.StatusBadRequest)
		return
	}

	err = rt.db.BanUser(ban.OwnerID, ban.PrayID)
	if err != nil {
		ctx.Logger.Info("Failed to ban user", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "text/plain")
	_, _ = w.Write([]byte("User banned!"))

	// ownerIdstr := ps.ByName("ownerID")

	// if ownerIdstr == "" {
	// 	http.Error(w, "missing ownerID", http.StatusBadRequest)
	// 	return
	// }

	// ownerID, err := strconv.Atoi(ownerIdstr)
	// if err != nil {
	// 	ctx.Logger.Info("Failed to convert OwnerID in int", err.Error())
	// 	http.Error(w, "ownerID not an int", http.StatusBadRequest)
	// 	return
	// }

	// prayIDStr := ps.ByName("prayID")

	// if prayIDStr == "" {
	// 	ctx.Logger.Info("Failed to get prayID", err.Error())
	// 	http.Error(w, "missing prayID", http.StatusBadRequest)
	// 	return
	// }

	// prayID, err := strconv.Atoi(prayIDStr)
	// if err != nil {
	// 	ctx.Logger.Info("Failed to convert prayID in int", err.Error())
	// 	http.Error(w, "prayID not an int", http.StatusBadRequest)
	// 	return
	// }

}
