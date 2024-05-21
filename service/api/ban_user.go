package api

import (
	"net/http"
	"strconv"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) BanUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	ownerIdstr := ps.ByName("ownerID")

	if ownerIdstr == "" {
		http.Error(w, "missing ownerID", http.StatusBadRequest)
		return
	}

	ownerID, err := strconv.Atoi(ownerIdstr)
	if err != nil {
		ctx.Logger.Info("Failed to convert OwnerID in int", err.Error())
		http.Error(w, "ownerID not an int", http.StatusBadRequest)
		return
	}

	prayIDStr := ps.ByName("prayID")

	if prayIDStr == "" {
		ctx.Logger.Info("Failed to get prayID", err.Error())
		http.Error(w, "missing prayID", http.StatusBadRequest)
		return
	}

	prayID, err := strconv.Atoi(prayIDStr)
	if err != nil {
		ctx.Logger.Info("Failed to convert prayID in int", err.Error())
		http.Error(w, "prayID not an int", http.StatusBadRequest)
		return
	}

	err = rt.db.BanUser(ownerID, prayID)
	if err != nil {
		ctx.Logger.Info("Failed to ban user", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "text/plain")
	_, _ = w.Write([]byte("User banned!"))

}
