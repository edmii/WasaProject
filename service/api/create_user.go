package api

import (
	"encoding/json"
	"net/http"

	"github.com/edmii/WasaProject/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

type User struct {
	Username string `json:"username"`
}

func (rt *_router) CreateUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	var user User

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

	err = rt.db.CreateUser(user.Username)
	if err != nil {
		ctx.Logger.Info("Failed to create user", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "text/plain")
	_, _ = w.Write([]byte("User created"))
}
