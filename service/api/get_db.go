package api

import (
	"encoding/json"
	"net/http"

	"github.com/edmii/WasaProject/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) getDB(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
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
		ctx.Logger.Info("Failed to get JSON encoding of DB", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "text/plain")
	_, _ = w.Write(encoded)
}
