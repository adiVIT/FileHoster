package files

import (
	"log"
	"time"
	"filestore/internal/db"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"os"
)

// StartFileDeletionJob starts a background job to delete expired files
func StartFileDeletionJob() {
	ticker := time.NewTicker(24 * time.Hour) // Run every 24 hours
	go func() {
		for range ticker.C {
			deleteExpiredFiles()
		}
	}()
}

// deleteExpiredFiles deletes files that are marked as expired
func deleteExpiredFiles() {
	// Initialize S3 session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")),
	})
	if err != nil {
		log.Printf("Failed to create AWS session: %v", err)
		return
	}

	s3Svc := s3.New(sess)

	// Query for expired files
	query := `SELECT id, s3_url FROM files WHERE expired = TRUE`
	rows, err := db.GetDB().Query(query)
	if err != nil {
		log.Printf("Failed to query expired files: %v", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var fileID uint
		var s3URL string
		if err := rows.Scan(&fileID, &s3URL); err != nil {
			log.Printf("Failed to scan file: %v", err)
			continue
		}

		// Delete file from S3
		_, err = s3Svc.DeleteObject(&s3.DeleteObjectInput{
			Bucket: aws.String(os.Getenv("S3_BUCKET")),
			Key:    aws.String(s3URL),
		})
		if err != nil {
			log.Printf("Failed to delete file from S3: %v", err)
			continue
		}

		// Delete file metadata from database
		_, err = db.GetDB().Exec(`DELETE FROM files WHERE id = $1`, fileID)
		if err != nil {
			log.Printf("Failed to delete file metadata: %v", err)
		}
	}
}