package files

import (
	"bytes"
	
	"fmt"
	
	
	"mime/multipart"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"filestore/internal/db"
)

// UploadFile handles the uploading of a file to S3 and saving its metadata
func UploadFile(userID uint, file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
	// Initialize S3 session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")),
	})
	if err != nil {
		return "", fmt.Errorf("failed to create AWS session: %w", err)
	}

	// Upload file to S3
	s3Svc := s3.New(sess)
	buffer := make([]byte, fileHeader.Size)
	file.Read(buffer)
	fileName := fmt.Sprintf("%d-%s", time.Now().Unix(), fileHeader.Filename)
	_, err = s3Svc.PutObject(&s3.PutObjectInput{
		Bucket:        aws.String(os.Getenv("S3_BUCKET")),
		Key:           aws.String(fileName),
		Body:          bytes.NewReader(buffer),
		ContentLength: aws.Int64(fileHeader.Size),
		ContentType:   aws.String(http.DetectContentType(buffer)),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file to S3: %w", err)
	}

	// Save file metadata in the database
	s3URL := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", os.Getenv("S3_BUCKET"), fileName)
	query := `INSERT INTO files (user_id, file_name, upload_date, size, s3_url) VALUES ($1, $2, $3, $4, $5)`
	_, err = db.GetDB().Exec(query, userID, fileHeader.Filename, time.Now(), fileHeader.Size, s3URL)
	if err != nil {
		return "", fmt.Errorf("failed to save file metadata: %w", err)
	}

	return s3URL, nil
}