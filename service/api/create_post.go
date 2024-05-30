package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

type Post struct {
	OwnerID int `json:"ownerID"`
}

func (rt *_router) CreatePost(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	if r.Method != "POST" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var post Post

	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		ctx.Logger.Info("Failed to decode request body ", err.Error())
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if post.OwnerID <= 0 {
		http.Error(w, "ownerID not valid", http.StatusBadRequest)
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
		ctx.Logger.Info("Failed to get cwd", err.Error())
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
		ctx.Logger.Info("Failed to save file", err.Error())
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "File %s uploaded successfully", header.Filename)

	err = rt.db.CreatePost(post.OwnerID, filepath)
	if err != nil {
		ctx.Logger.Info("Failed to create post", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "text/plain")
	_, _ = w.Write([]byte("Post created"))
}
