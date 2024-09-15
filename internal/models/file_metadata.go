// internal/models/file_metadata.go

package models
import (
	"time"
)
type FileMetadata struct {
	ID        uint
	UserID    uint
	FileName  string
	UploadDate time.Time
	Size      int64
	S3URL     string
}
