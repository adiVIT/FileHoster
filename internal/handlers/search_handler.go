package handlers

import (
    "net/http"
    "time"
    "filestore/internal/files"
    "filestore/internal/utils"
    "filestore/internal/auth"
)

func SearchFilesHandler(w http.ResponseWriter, r *http.Request) {
    userID, err := auth.ValidateJWT(r.Header.Get("Authorization"))
    if err != nil {
        utils.RespondWithError(w, http.StatusUnauthorized, "Invalid token")
        return
    }

    params := files.SearchParams{
        FileName: r.URL.Query().Get("name"),
        FileType: r.URL.Query().Get("type"),
    }

    if dateStr := r.URL.Query().Get("date"); dateStr != "" {
        params.UploadDate, _ = time.Parse("2006-01-02", dateStr)
    }

    results, err := files.SearchFiles(userID, params)
    if err != nil {
        utils.RespondWithError(w, http.StatusInternalServerError, "Error searching files")
        return
    }

    utils.RespondWithJSON(w, http.StatusOK, results)
}