package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) DestroyDB(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	err := rt.db.DestroyDB()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "text/plain")
	_, _ = w.Write([]byte("Database destroyed"))
}
