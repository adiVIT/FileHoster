package handlers

import (
	"net/http"
	"github.com/gorilla/mux"
	"filestore/internal/auth"
	"filestore/internal/files"
	"filestore/internal/utils"
	"strconv"
)

// ShareFileHandler handles requests to share a file via a public link
func ShareFileHandler(w http.ResponseWriter, r *http.Request) {
	// Authenticate user
	_, err := auth.ValidateJWT(r.Header.Get("Authorization"))
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Invalid token")
		return
	}

	// Get file ID from URL parameters
	vars := mux.Vars(r)
	fileIDStr := vars["file_id"]
	fileID, err := strconv.ParseUint(fileIDStr, 10, 32)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid file ID")
		return
	}

	// Retrieve the public URL for the file
	publicURL, err := files.ShareFile(uint(fileID))
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error sharing file")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"url": publicURL})
}