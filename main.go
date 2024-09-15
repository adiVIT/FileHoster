package main

import (
    "log"
    "net/http"
    "os"

    "github.com/gorilla/mux"
    "github.com/joho/godotenv"
    "filestore/internal/handlers"
    "filestore/internal/db"
    "filestore/internal/cache"
    "filestore/internal/files" // Import the files package
)

func main() {
    // Load environment variables
    err := godotenv.Load()
    if err != nil {
        log.Fatalf("Error loading .env file")
    }

    // Initialize database and cache
    db.InitDB()
    cache.InitRedis()
    db.CreateFileTable()
    db.CreateUserTable()
    db.CreateFileTable()

    // Start the background job for file deletion
    files.StartFileDeletionJob()

    // Create a new router
    r := mux.NewRouter()

    // Register handlers
    r.HandleFunc("/register", handlers.RegisterHandler).Methods("POST")
    r.HandleFunc("/login", handlers.LoginHandler).Methods("POST")
    r.HandleFunc("/upload", handlers.UploadFileHandler).Methods("POST")
    r.HandleFunc("/files", handlers.ListFilesHandler).Methods("GET")
    r.HandleFunc("/share/{file_id}", handlers.ShareFileHandler).Methods("GET")
    r.HandleFunc("/search", handlers.SearchFilesHandler).Methods("GET")

    // Start the server
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080" // Default port if not specified
    }
    log.Printf("Server starting on port %s", port)
    log.Fatal(http.ListenAndServe(":"+port, r))
}