package api

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

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

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current working directory:", err)
		return
	}

	// Define your relative path
	relativePath := "service/database/images-db"

	// Join the current working directory with the relative path
	uploadPath := filepath.Join(cwd, relativePath)

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
