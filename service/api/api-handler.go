package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes
	rt.router.GET("/", rt.getHelloWorld)
	rt.router.GET("/context", rt.wrap(rt.getContextReply))
	rt.router.GET("/feed", rt.getFeed)
	rt.router.GET("/db/:table", rt.getDB)
	rt.router.GET("/createuser/:username", rt.CreateUser)
	rt.router.GET("/DESTROYDB/sure", rt.DestroyDB)

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}

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

func (rt *_router) DestroyDB(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	err := rt.db.DestroyDB()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "text/plain")
	_, _ = w.Write([]byte("Database destroyed"))
}
