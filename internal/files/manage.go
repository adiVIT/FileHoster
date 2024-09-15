package files

import (
	"database/sql"
	"fmt"
	

	"filestore/internal/db"
	"filestore/internal/cache"
	"filestore/internal/models"
)

// models.FileMetadata represents the metadata of a file


// ListFiles retrieves all files for a specific user
func ListFiles(userID uint) ([]models.FileMetadata, error) {
	// Check cache first
	cachedFiles, found := cache.GetFilesFromCache(userID)
	if found {
		return cachedFiles, nil
	}

	// Query database if not found in cache
	query := `SELECT id, user_id, file_name, upload_date, size, s3_url FROM files WHERE user_id = $1`
	rows, err := db.GetDB().Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query files: %w", err)
	}
	defer rows.Close()

	var files []models.FileMetadata
	for rows.Next() {
		var file models.FileMetadata
		if err := rows.Scan(&file.ID, &file.UserID, &file.FileName, &file.UploadDate, &file.Size, &file.S3URL); err != nil {
			return nil, fmt.Errorf("failed to scan file: %w", err)
		}
		files = append(files, file)
	}

	// Cache the result
	cache.SetFilesInCache(userID, files)

	return files, nil
}

// ShareFile generates a public URL for sharing a file
func ShareFile(fileID uint) (string, error) {
	query := `SELECT s3_url FROM files WHERE id = $1`
	var s3URL string
	err := db.GetDB().QueryRow(query, fileID).Scan(&s3URL)
	if err == sql.ErrNoRows {
		return "", fmt.Errorf("file not found")
	} else if err != nil {
		return "", fmt.Errorf("failed to query file: %w", err)
	}

	// Generate a public link (could include additional logic for URL expiration)
	return s3URL, nil
}