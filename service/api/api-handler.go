package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
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
	rt.router.POST("/createpost/:ownerID", rt.wrap(rt.CreatePost))

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

func (rt *_router) CreatePost(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	if r.Method != "POST" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	ownerIDStr := ps.ByName("ownerID")

	if ownerIDStr == "" {
		http.Error(w, "missing ownerID", http.StatusBadRequest)
		return
	}

	ownerID, err := strconv.Atoi(ownerIDStr)
	if err != nil {
		http.Error(w, "ownerID not an int", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		ctx.Logger.Info("Failed to read file from request", err.Error())
		http.Error(w, "Failed to read file from request", http.StatusBadRequest)
		return
	}
	defer file.Close()

	uploadPath := "service/database/images-db"

	filepath := filepath.Join(uploadPath, header.Filename)

	out, err := os.Create(filepath)
	if err != nil {
		ctx.Logger.Info("Failed create file from request", err.Error())
		http.Error(w, "Failed to create file", http.StatusInternalServerError)
		return
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "File %s uploaded successfully", header.Filename)

	err = rt.db.CreatePost(ownerID, filepath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "text/plain")
	_, _ = w.Write([]byte("Post created"))
}
