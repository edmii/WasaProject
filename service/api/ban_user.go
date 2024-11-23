package api

import (
	"encoding/json"
	"net/http"

	// "strconv"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
	"github.com\edmii\WasaProject\service\models\structs.go"
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

	result, err := rt.db.BanUser(ban.OwnerID, ban.PrayID)
	if err != nil {
		ctx.Logger.Info("Failed to ban user", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "text/plain")
	switch result {
	case 1:
		_, _ = w.Write([]byte("User unbanned!"))
	case 2:
		_, _ = w.Write([]byte("User banned!"))
	}

}

func (rt *_router) GetBans(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
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

	// Get the list of banned users from the database for the given OwnerID
	bannedUsers, err := rt.db.GetBannedUsers(ban.OwnerID)
	if err != nil {
		ctx.Logger.Info("Failed to retrieve banned users", err.Error())
		http.Error(w, "Failed to retrieve banned users", http.StatusInternalServerError)
		return
	}

	response := map[string][]int{
		"prayID": bannedUsers,
	}

	// Convert the list to JSON and send the response
	w.Header().Set("content-type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		ctx.Logger.Info("Failed to encode banned users to JSON", err.Error())
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
