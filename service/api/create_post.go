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

	"github.com/edmii/WasaProject/service/api/reqcontext"
	structs "github.com/edmii/WasaProject/service/models"
	"github.com/edmii/WasaProject/service/utils"
	"github.com/julienschmidt/httprouter"
)

// type Post struct {
// 	PostID    int    `json:"postID"`
// 	OwnerID   int    `json:"ownerID"`
// 	Directory string `json:"imagePath"`
// 	PostedAt  string `json:"postedAt"`

// 	RequesterID int `json:"requesterID"`
// }

func (rt *_router) uploadPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	if r.Method != "POST" {
		utils.SendErrorResponse(w, "Invalid request method", []string{"Only POST requests are allowed"}, http.StatusMethodNotAllowed)
		return
	}

	// Parse the multipart form with a size limit (e.g., 10 MB)
	err := r.ParseMultipartForm(10 << 20) // 10MB
	if err != nil {
		ctx.Logger.Info("Failed to parse multipart form ", err.Error())
		utils.SendErrorResponse(w, "Invalid multipart form", []string{"Failed to parse multipart form"}, http.StatusBadRequest)
		return
	}

	// Extract ownerID from the form data
	ownerIDStr := r.FormValue("ownerID")
	if ownerIDStr == "" {
		utils.SendErrorResponse(w, "Invalid request", []string{"ownerID is required"}, http.StatusBadRequest)
		return
	}

	var ownerID int
	_, err = fmt.Sscanf(ownerIDStr, "%d", &ownerID)
	if err != nil || ownerID <= 0 {
		utils.SendErrorResponse(w, "Invalid request", []string{"Invalid ownerID value"}, http.StatusBadRequest)
		return
	}

	// Extract the file from the form
	file, header, err := r.FormFile("file")
	if err != nil {
		ctx.Logger.Info("Failed to read file from request", err.Error())
		utils.SendErrorResponse(w, "Invalid file", []string{"Failed to read file from request"}, http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Get the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		ctx.Logger.Info("Failed to get cwd", err.Error())
		utils.SendErrorResponse(w, "Server error", []string{"Failed to get current working directory", err.Error()}, http.StatusInternalServerError)
		return
	}

	relativePath := "service/database/images-db"
	uploadPath := filepath.Join(cwd, relativePath)

	PostedAt := time.Now()

	// Use the extracted data to create a post
	PostID, err := rt.db.CreatePost(ownerID, uploadPath, PostedAt)
	if err != nil {
		ctx.Logger.Info("Failed to create post", err.Error())
		utils.SendErrorResponse(w, "Database error", []string{"Failed to create post", err.Error()}, http.StatusInternalServerError)
		return
	}

	postIDStr := strconv.Itoa(PostID)
	filepath := filepath.Join(uploadPath, postIDStr)

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		ctx.Logger.Info("Failed to create file from request", err.Error())
		utils.SendErrorResponse(w, "Server error", []string{"Failed to create file", err.Error()}, http.StatusInternalServerError)

		_ = rt.db.DeletePost(PostID, -1)
		return
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		ctx.Logger.Info("Failed to save file", err.Error())
		utils.SendErrorResponse(w, "Server error", []string{"Failed to save file", err.Error()}, http.StatusInternalServerError)

		_ = rt.db.DeletePost(PostID, -1)
		_ = os.Remove(filepath)
		return
	}

	ctx.Logger.Info(fmt.Sprintf("File %s uploaded successfully", header.Filename))

	w.Header().Set("content-type", "text/plain")
	_, _ = w.Write([]byte("Post created"))
}

// func (rt *_router) CreatePost(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

// 	if r.Method != "POST" {
// 		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
// 		return
// 	}

// 	// Parse the multipart form with a size limit (e.g., 10 MB)
// 	err := r.ParseMultipartForm(10 << 20) // 10MB
// 	if err != nil {
// 		ctx.Logger.Info("Failed to parse multipart form ", err.Error())
// 		http.Error(w, "Failed to parse multipart form", http.StatusBadRequest)
// 		return
// 	}

// 	// Extract ownerID from the form data
// 	ownerIDStr := r.FormValue("ownerID")
// 	if ownerIDStr == "" {
// 		http.Error(w, "ownerID is required", http.StatusBadRequest)
// 		return
// 	}

// 	var ownerID int
// 	_, err = fmt.Sscanf(ownerIDStr, "%d", &ownerID)
// 	if err != nil || ownerID <= 0 {
// 		http.Error(w, "Invalid ownerID value", http.StatusBadRequest)
// 		return
// 	}

// 	// Extract the file from the form
// 	file, header, err := r.FormFile("file")
// 	if err != nil {
// 		ctx.Logger.Info("Failed to read file from request", err.Error())
// 		http.Error(w, "Failed to read file from request", http.StatusBadRequest)
// 		return
// 	}
// 	defer file.Close()

// 	// Get the current working directory
// 	cwd, err := os.Getwd()
// 	if err != nil {
// 		ctx.Logger.Info("Failed to get cwd", err.Error())
// 		fmt.Println("Error getting current working directory:", err)
// 		return
// 	}

// 	// Define your relative path for saving the file
// 	relativePath := "service/database/images-db"

// 	// Join the current working directory with the relative path
// 	uploadPath := filepath.Join(cwd, relativePath)

// 	PostedAt := time.Now()

// 	// Use the extracted data to create a post
// 	PostID, err := rt.db.CreatePost(ownerID, uploadPath, PostedAt)
// 	if err != nil {
// 		ctx.Logger.Info("Failed to create post", err.Error())
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	postIDStr := strconv.Itoa(PostID)
// 	// Create the full filepath for saving the file
// 	filepath := filepath.Join(uploadPath, postIDStr)

// 	// Create the file
// 	out, err := os.Create(filepath)
// 	if err != nil {
// 		ctx.Logger.Info("Failed create file from request", err.Error())
// 		http.Error(w, "Failed to create file", http.StatusInternalServerError)

// 		err = rt.db.DeletePost(PostID, -1)
// 		if err != nil {
// 			ctx.Logger.Info("Failed to delete post", err.Error())
// 			http.Error(w, "Failed to delete post", http.StatusInternalServerError)
// 		}

// 		return
// 	}
// 	defer out.Close()

// 	// Save the file content to the specified location
// 	_, err = io.Copy(out, file)
// 	if err != nil {
// 		ctx.Logger.Info("Failed to save file", err.Error())
// 		http.Error(w, "Failed to save file", http.StatusInternalServerError)
// 		err = rt.db.DeletePost(PostID, -1)
// 		if err != nil {
// 			ctx.Logger.Info("Failed to delete post", err.Error())
// 			http.Error(w, "Failed to delete post", http.StatusInternalServerError)
// 		}

// 		//delete filepath
// 		err = os.Remove(filepath)
// 		if err != nil {
// 			ctx.Logger.Info("Failed to delete file", err.Error())
// 			http.Error(w, "Failed to delete file", http.StatusInternalServerError)
// 		}

// 		return
// 	}

// 	ctx.Logger.Info(fmt.Sprintf("File %s uploaded successfully", header.Filename))

// 	w.Header().Set("content-type", "text/plain")
// 	_, _ = w.Write([]byte("Post created"))
// }

func (rt *_router) deletePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var post structs.Post
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		ctx.Logger.Info("Failed to decode request body ", err.Error())
		utils.SendErrorResponse(w, "Invalid request body", []string{"Failed to decode JSON request body"}, http.StatusBadRequest)
		return
	}

	if post.RequesterID <= 0 {
		utils.SendErrorResponse(w, "Invalid request", []string{"Invalid requester ID"}, http.StatusBadRequest)
		return
	}

	err = rt.db.DeletePost(post.PostID, post.RequesterID)
	if err != nil {
		ctx.Logger.Info("Failed to delete post", err.Error())
		utils.SendErrorResponse(w, "Database error", []string{"Failed to delete post", err.Error()}, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	_, _ = w.Write([]byte("Post deleted"))
}

// func (rt *_router) DeletePost(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
// 	var post structs.Post
// 	err := json.NewDecoder(r.Body).Decode(&post)
// 	if err != nil {
// 		ctx.Logger.Info("Failed to decode request body ", err.Error())
// 		http.Error(w, "Invalid request body", http.StatusBadRequest)
// 		return
// 	}

// 	if post.RequesterID <= 0 {
// 		http.Error(w, "invalid requester ID", http.StatusBadRequest)
// 		return
// 	}

// 	err = rt.db.DeletePost(post.PostID, post.RequesterID)
// 	if err != nil {
// 		ctx.Logger.Info("Failed to delete post", err.Error())
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("content-type", "text/plain")
// 	_, _ = w.Write([]byte("Post deleted"))
// }

func (rt *_router) getUserPosts(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var user structs.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		ctx.Logger.Info("Failed to decode request body ", err.Error())
		utils.SendErrorResponse(w, "Invalid request body", []string{"Failed to decode JSON request body"}, http.StatusBadRequest)
		return
	}

	if user.Username == "" {
		utils.SendErrorResponse(w, "Invalid request", []string{"Username is required"}, http.StatusBadRequest)
		return
	}

	posts, err := rt.db.GetUserPosts(user.Username)
	if err != nil {
		ctx.Logger.Info("Failed to get posts", err.Error())
		utils.SendErrorResponse(w, "Database error", []string{"Failed to retrieve user posts", err.Error()}, http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"status":  "success",
		"message": "User posts retrieved",
		"data":    posts,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		ctx.Logger.Info("Failed to encode posts", err.Error())
		utils.SendErrorResponse(w, "Server error", []string{"Failed to encode response", err.Error()}, http.StatusInternalServerError)
	}
}

// func (rt *_router) GetUserPosts(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
// 	var user structs.User

// 	err := json.NewDecoder(r.Body).Decode(&user)
// 	if err != nil {
// 		ctx.Logger.Info("Failed to decode request body ", err.Error())
// 		http.Error(w, "Invalid request body", http.StatusBadRequest)
// 		return
// 	}

// 	posts, err := rt.db.GetUserPosts(user.Username)
// 	if err != nil {
// 		ctx.Logger.Info("Failed to get posts", err.Error())
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("content-type", "application/json")

// 	response := map[string]interface{}{
// 		"status":  "success",
// 		"message": "User posts retrieved",
// 		"data":    posts,
// 	}

// 	err = json.NewEncoder(w).Encode(response)
// 	if err != nil {
// 		ctx.Logger.Info("Failed to encode posts", err.Error())
// 		http.Error(w, "Failed to encode posts", http.StatusInternalServerError)
// 		return
// 	}
// }
