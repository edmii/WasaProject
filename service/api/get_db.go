package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) getDB(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	table := ps.ByName("table")
	if table == "" {
		http.Error(w, "missing table name", http.StatusBadRequest)
		return
	}

	content, err := rt.db.GetDatabaseTableContent(table)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	encoded, err := json.Marshal(content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "text/plain")
	_, _ = w.Write(encoded)
}
