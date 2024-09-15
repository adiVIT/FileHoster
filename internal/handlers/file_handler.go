package handlers

import (
	// "encoding/json"
	"net/http"
	"filestore/internal/auth"
	"filestore/internal/files"
	"filestore/internal/utils"
	"log"
	// "mime/multipart"
	// "strconv"
)

// UploadFileHandler handles file upload requests
func UploadFileHandler(w http.ResponseWriter, r *http.Request) {
	// Authenticate user
	userID, err := auth.ValidateJWT(r.Header.Get("Authorization"))
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Invalid token")
		return
	}

	// Parse the multipart form
	err = r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Unable to parse form")
		return
	}

	// Retrieve the file from form data
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid file")
		return
	}
	defer file.Close()

	// Upload the file
	s3URL, err := files.UploadFile(userID, file, fileHeader)
	if err != nil {
		log.Printf("Error creating user: %v", err)

		utils.RespondWithError(w, http.StatusInternalServerError, "Error uploading file")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"url": s3URL})
}

// ListFilesHandler handles requests to list user files
func ListFilesHandler(w http.ResponseWriter, r *http.Request) {
	// Authenticate user
	userID, err := auth.ValidateJWT(r.Header.Get("Authorization"))
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Invalid token")
		return
	}

	// Retrieve files for the user
	files, err := files.ListFiles(userID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error retrieving files")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, files)
}