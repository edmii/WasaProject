package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) CreateUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	username := ps.ByName("username")
	if username == "" {
		http.Error(w, "missing username", http.StatusBadRequest)
		return
	}

	err := rt.db.CreateUser(username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "text/plain")
	_, _ = w.Write([]byte("User created"))
}
