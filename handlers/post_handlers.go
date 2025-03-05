package handlers

import (
    "blog-api/models"
    "blog-api/services"
    "encoding/json"
    "net/http"
    "strconv"

    "github.com/gorilla/mux"
)

type PostHandler struct {
    PostService *services.PostService
}

func (h *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request){
	var post models.Post
    // Read from the request body, decode and put it into post
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

	id, err := h.PostService.CreatePost(&post)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

	post.ID = id
	w.WriteHeader(http.StatusCreated) // set status to 201
	json.NewEncoder(w).Encode(post) // encode the post back and send it to the writer
}

func (h *PostHandler) FetchPost(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil || id <= 0 {
        http.Error(w, "Invalid post ID", http.StatusBadRequest)
        return
    }

	post, err := h.PostService.FetchPost(id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

	if post == nil {
        http.Error(w, "Post not found", http.StatusNotFound)
        return
    }
	json.NewEncoder(w).Encode(post)
}

func (h *PostHandler) FetchPosts(w http.ResponseWriter, r *http.Request) {
    term := r.URL.Query().Get("search") // basically gets something from search=something

    posts, err := h.PostService.FetchPosts(term)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(posts)
}

func (h *PostHandler) UpdatePost(w http.ResponseWriter, r *http.Request) {
	var post models.Post
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil || id <= 0{
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }
    post.ID = id
	err = h.PostService.UpdatePost(&post)
	if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }


	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Post updated successfully"})
}

func (h *PostHandler) DeletePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil || id <= 0{
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	err = h.PostService.DeletePost(id)
	if err != nil {
        if err.Error() == "post not found" {
            http.Error(w, "Post not found", http.StatusNotFound)
        } else {
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
        return
    }

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Post deleted successfully"})
}