package api

import (
	"encoding/json"
	"net/http"

	// "strconv"

	"github.com/edmii/WasaProject/service/api/reqcontext"
	structs "github.com/edmii/WasaProject/service/models"
	"github.com/edmii/WasaProject/service/utils"
	"github.com/julienschmidt/httprouter"
)

// type Ban struct {
// 	OwnerID int `json:"ownerID"`
// 	PrayID  int `json:"prayID"`
// }

func (rt *_router) banUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var ban structs.Ban

	err := json.NewDecoder(r.Body).Decode(&ban)
	if err != nil {
		ctx.Logger.Info("Failed to decode request body ", err.Error())
		utils.SendErrorResponse(w, "Invalid request body", []string{"Failed to decode JSON request body"}, http.StatusBadRequest)
		return
	}

	if ban.OwnerID <= 0 {
		utils.SendErrorResponse(w, "Invalid request", []string{"Invalid OwnerID"}, http.StatusBadRequest)
		return
	}

	if ban.PrayID <= 0 {
		utils.SendErrorResponse(w, "Invalid request", []string{"Invalid PrayID"}, http.StatusBadRequest)
		return
	}

	result, err := rt.db.BanUser(ban.OwnerID, ban.PrayID)
	if err != nil {
		ctx.Logger.Info("Failed to ban user", err.Error())
		utils.SendErrorResponse(w, "Database error", []string{"Failed to ban/unban user", err.Error()}, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{}

	switch result {
	case 1:
		response["message"] = "User already banned!"
	case 2:
		response["message"] = "User banned!"
	default:
		utils.SendErrorResponse(w, "Unexpected result", []string{"Unknown result code from BanUser operation"}, http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		ctx.Logger.Info("Failed to encode response", err.Error())
		utils.SendErrorResponse(w, "Server error", []string{"Failed to encode response", err.Error()}, http.StatusInternalServerError)
	}
}

func (rt *_router) unbanUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var ban structs.Ban

	err := json.NewDecoder(r.Body).Decode(&ban)
	if err != nil {
		ctx.Logger.Info("Failed to decode request body ", err.Error())
		utils.SendErrorResponse(w, "Invalid request body", []string{"Failed to decode JSON request body"}, http.StatusBadRequest)
		return
	}

	if ban.OwnerID <= 0 {
		utils.SendErrorResponse(w, "Invalid request", []string{"Invalid OwnerID"}, http.StatusBadRequest)
		return
	}

	if ban.PrayID <= 0 {
		utils.SendErrorResponse(w, "Invalid request", []string{"Invalid PrayID"}, http.StatusBadRequest)
		return
	}

	result, err := rt.db.BanUser(ban.OwnerID, ban.PrayID)
	if err != nil {
		ctx.Logger.Info("Failed to ban user", err.Error())
		utils.SendErrorResponse(w, "Database error", []string{"Failed to ban/unban user", err.Error()}, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{}

	switch result {
	case 1:
		response["message"] = "User unbanned!"
	case 2:
		response["message"] = "User already not banned!"
	default:
		utils.SendErrorResponse(w, "Unexpected result", []string{"Unknown result code from BanUser operation"}, http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		ctx.Logger.Info("Failed to encode response", err.Error())
		utils.SendErrorResponse(w, "Server error", []string{"Failed to encode response", err.Error()}, http.StatusInternalServerError)
	}
}

// func (rt *_router) BanUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

// 	var ban structs.Ban

// 	err := json.NewDecoder(r.Body).Decode(&ban)
// 	if err != nil {
// 		ctx.Logger.Info("Failed to decode request body ", err.Error())
// 		http.Error(w, "Invalid request body", http.StatusBadRequest)
// 		return
// 	}

// 	if ban.OwnerID <= 0 {
// 		http.Error(w, "Invalid OwnerID", http.StatusBadRequest)
// 		return
// 	}

// 	if ban.PrayID <= 0 {
// 		http.Error(w, "Invalid PrayID", http.StatusBadRequest)
// 		return
// 	}

// 	result, err := rt.db.BanUser(ban.OwnerID, ban.PrayID)
// 	if err != nil {
// 		ctx.Logger.Info("Failed to ban user", err.Error())
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("content-type", "text/plain")
// 	switch result {
// 	case 1:
// 		_, _ = w.Write([]byte("User unbanned!"))
// 	case 2:
// 		_, _ = w.Write([]byte("User banned!"))
// 	}

// }

func (rt *_router) GetBans(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var ban structs.Ban

	err := json.NewDecoder(r.Body).Decode(&ban)
	if err != nil {
		ctx.Logger.Info("Failed to decode request body ", err.Error())
		utils.SendErrorResponse(w, "Invalid request body", []string{"Failed to decode JSON request body"}, http.StatusBadRequest)
		return
	}

	if ban.OwnerID <= 0 {
		utils.SendErrorResponse(w, "Invalid request", []string{"Invalid OwnerID"}, http.StatusBadRequest)
		return
	}

	// Get the list of banned users from the database for the given OwnerID
	bannedUsers, err := rt.db.GetBannedUsers(ban.OwnerID)
	if err != nil {
		ctx.Logger.Info("Failed to retrieve banned users", err.Error())
		utils.SendErrorResponse(w, "Database error", []string{"Failed to retrieve banned users", err.Error()}, http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"status":  "success",
		"message": "Banned users retrieved",
		"data":    bannedUsers,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		ctx.Logger.Info("Failed to encode banned users to JSON", err.Error())
		utils.SendErrorResponse(w, "Server error", []string{"Failed to encode response", err.Error()}, http.StatusInternalServerError)
	}
}

// func (rt *_router) GetBans(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
// 	var ban structs.Ban

// 	err := json.NewDecoder(r.Body).Decode(&ban)
// 	if err != nil {
// 		ctx.Logger.Info("Failed to decode request body ", err.Error())
// 		http.Error(w, "Invalid request body", http.StatusBadRequest)
// 		return
// 	}

// 	if ban.OwnerID <= 0 {
// 		http.Error(w, "Invalid OwnerID", http.StatusBadRequest)
// 		return
// 	}

// 	// Get the list of banned users from the database for the given OwnerID
// 	bannedUsers, err := rt.db.GetBannedUsers(ban.OwnerID)
// 	if err != nil {
// 		ctx.Logger.Info("Failed to retrieve banned users", err.Error())
// 		http.Error(w, "Failed to retrieve banned users", http.StatusInternalServerError)
// 		return
// 	}

// 	response := map[string]interface{}{
// 		"status":  "success",
// 		"message": "Banned users retrieved",
// 		"data":    bannedUsers,
// 	}

// 	// Convert the list to JSON and send the response
// 	w.Header().Set("content-type", "application/json")
// 	err = json.NewEncoder(w).Encode(response)
// 	if err != nil {
// 		ctx.Logger.Info("Failed to encode banned users to JSON", err.Error())
// 		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
// 		return
// 	}
// }
