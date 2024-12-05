package api

import (
	"encoding/json"
	"net/http"

	"github.com/edmii/WasaProject/service/api/reqcontext"
	structs "github.com/edmii/WasaProject/service/models"
	"github.com/julienschmidt/httprouter"
)

// type User struct {
// 	Username string `json:"username"`
// }

func (rt *_router) Login(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var user structs.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		ctx.Logger.Info("Failed to decode request body ", err.Error())
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if user.Username == "" {
		http.Error(w, "missing username", http.StatusBadRequest)
		return
	}

	var exist bool
	exist, err = rt.db.CheckUserExist(user.Username)
	if err != nil {
		ctx.Logger.Info("Failed to check user existance", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !exist {
		err = rt.db.CreateUser(user.Username)
		if err != nil {
			ctx.Logger.Info("Failed to create user", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	Map := map[string]interface{}{
		"msg":  "User logged in",
		"user": user,
	}

	json, err := json.Marshal(Map)
	if err != nil {
		ctx.Logger.Info("Failed to marshal json", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "text/plain")
	_, _ = w.Write(json)
}

func (rt *_router) ChangeUsername(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var user structs.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		ctx.Logger.Info("Failed to decode request body ", err.Error())
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	newUsername := user.Username

	err = rt.db.ChangeUsername(user.ID, newUsername)

	if err != nil {
		ctx.Logger.Info("Failed to create user", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "text/plain")
	_, _ = w.Write([]byte("Username changed"))
}
