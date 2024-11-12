package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

type Post struct {
	PostID    int    `json:"postID"`
	OwnerID   int    `json:"ownerID"`
	Directory string `json:"imagePath"`
	PostedAt  string `json:"postedAt"`

	RequesterID int `json:"requesterID"`
}

func (rt *_router) CreatePost(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	if r.Method != "POST" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Parse the multipart form with a size limit (e.g., 10 MB)
	err := r.ParseMultipartForm(10 << 20) // 10MB
	if err != nil {
		ctx.Logger.Info("Failed to parse multipart form ", err.Error())
		http.Error(w, "Failed to parse multipart form", http.StatusBadRequest)
		return
	}

	// Extract ownerID from the form data
	ownerIDStr := r.FormValue("ownerID")
	if ownerIDStr == "" {
		http.Error(w, "ownerID is required", http.StatusBadRequest)
		return
	}

	var ownerID int
	_, err = fmt.Sscanf(ownerIDStr, "%d", &ownerID)
	if err != nil || ownerID <= 0 {
		http.Error(w, "Invalid ownerID value", http.StatusBadRequest)
		return
	}

	// Extract the file from the form
	file, header, err := r.FormFile("file")
	if err != nil {
		ctx.Logger.Info("Failed to read file from request", err.Error())
		http.Error(w, "Failed to read file from request", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Get the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		ctx.Logger.Info("Failed to get cwd", err.Error())
		fmt.Println("Error getting current working directory:", err)
		return
	}

	// Define your relative path for saving the file
	relativePath := "service/database/images-db"

	// Join the current working directory with the relative path
	uploadPath := filepath.Join(cwd, relativePath)

	PostedAt := time.Now()

	// Use the extracted data to create a post
	PostID, err := rt.db.CreatePost(ownerID, uploadPath, PostedAt)
	if err != nil {
		ctx.Logger.Info("Failed to create post", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	postIDStr := strconv.Itoa(PostID)
	// Create the full filepath for saving the file
	filepath := filepath.Join(uploadPath, postIDStr)

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		ctx.Logger.Info("Failed create file from request", err.Error())
		http.Error(w, "Failed to create file", http.StatusInternalServerError)

		err = rt.db.DeletePost(PostID, -1)
		if err != nil {
			ctx.Logger.Info("Failed to delete post", err.Error())
			http.Error(w, "Failed to delete post", http.StatusInternalServerError)
		}

		return
	}
	defer out.Close()

	// Save the file content to the specified location
	_, err = io.Copy(out, file)
	if err != nil {
		ctx.Logger.Info("Failed to save file", err.Error())
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		err = rt.db.DeletePost(PostID, -1)
		if err != nil {
			ctx.Logger.Info("Failed to delete post", err.Error())
			http.Error(w, "Failed to delete post", http.StatusInternalServerError)
		}

		//delete filepath
		err = os.Remove(filepath)
		if err != nil {
			ctx.Logger.Info("Failed to delete file", err.Error())
			http.Error(w, "Failed to delete file", http.StatusInternalServerError)
		}

		return
	}

	ctx.Logger.Info(fmt.Sprintf("File %s uploaded successfully", header.Filename))

	w.Header().Set("content-type", "text/plain")
	_, _ = w.Write([]byte("Post created"))
}

func (rt *_router) DeletePost(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var post Post
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		ctx.Logger.Info("Failed to decode request body ", err.Error())
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if post.RequesterID <= 0 {
		http.Error(w, "invalid requester ID", http.StatusBadRequest)
		return
	}

	err = rt.db.DeletePost(post.PostID, post.RequesterID)
	if err != nil {
		ctx.Logger.Info("Failed to delete post", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "text/plain")
	_, _ = w.Write([]byte("Post deleted"))
}

func (rt *_router) GetUserPosts(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var user User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		ctx.Logger.Info("Failed to decode request body ", err.Error())
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	posts, err := rt.db.GetUserPosts(user.Username)
	if err != nil {
		ctx.Logger.Info("Failed to get posts", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "application/json")
	err = json.NewEncoder(w).Encode(posts)
	if err != nil {
		ctx.Logger.Info("Failed to encode posts", err.Error())
		http.Error(w, "Failed to encode posts", http.StatusInternalServerError)
		return
	}

}
