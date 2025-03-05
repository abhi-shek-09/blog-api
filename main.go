package main

import (
	"blog-api/database"
	"blog-api/handlers"
	"blog-api/repository"
	"github.com/gorilla/mux"
	"blog-api/services"
	"log"
	"net/http"
)

func main() {
    database.ConnectDB() 
    defer database.CloseDB()

    repo := &repository.PostRepository{DB: database.DB}
    service := services.NewPostService(repo)
    handler := &handlers.PostHandler{PostService: service}

    r := mux.NewRouter()

	r.HandleFunc("/post", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello from /post!"))
	}).Methods("GET")
    r.HandleFunc("/post", handler.CreatePost).Methods("POST")
    r.HandleFunc("/post/{id:[0-9]+}", handler.FetchPost).Methods("GET")
    r.HandleFunc("/posts", handler.FetchPosts).Methods("GET")
    r.HandleFunc("/post/{id:[0-9]+}", handler.UpdatePost).Methods("PUT")
    r.HandleFunc("/post/{id:[0-9]+}", handler.DeletePost).Methods("DELETE")

    log.Println("Server started on :8080")
    log.Fatal(http.ListenAndServe(":8080", r))
}