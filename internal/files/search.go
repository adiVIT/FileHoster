package files

import (
    "fmt"
    "strings"
    "time"

    "filestore/internal/db"
    "filestore/internal/models"
)

type SearchParams struct {
    FileName   string
    UploadDate time.Time
    FileType   string
}

func SearchFiles(userID uint, params SearchParams) ([]models.FileMetadata, error) {
    query := `SELECT id, user_id, file_name, upload_date, size, s3_url FROM files WHERE user_id = $1`
    args := []interface{}{userID}

    if params.FileName != "" {
        query += " AND file_name ILIKE $" + fmt.Sprint(len(args)+1)
        args = append(args, "%"+params.FileName+"%")
    }

    if !params.UploadDate.IsZero() {
        query += " AND upload_date::date = $" + fmt.Sprint(len(args)+1)
        args = append(args, params.UploadDate)
    }

    if params.FileType != "" {
        query += " AND file_name ILIKE $" + fmt.Sprint(len(args)+1)
        args = append(args, "%."+strings.ToLower(params.FileType))
    }

    rows, err := db.GetDB().Query(query, args...)
    if err != nil {
        return nil, fmt.Errorf("failed to search files: %w", err)
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

    return files, nil
}