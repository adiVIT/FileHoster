package db

import (
	"time"
)

// User represents a user in the system
type User struct {
	ID       uint   `json:"id"`
	Email    string `json:"email"`
	Password string `json:"-"`
}

// File represents a file uploaded by a user
type File struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"user_id"`
	FileName  string    `json:"file_name"`
	UploadDate time.Time `json:"upload_date"`
	Size      int64     `json:"size"`
	S3URL     string    `json:"s3_url"`
}

// CreateUserTable creates the users table in the database
func CreateUserTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		email VARCHAR(255) UNIQUE NOT NULL,
		password VARCHAR(255) NOT NULL
	)`
	_, err := GetDB().Exec(query)
	return err
}

// CreateFileTable creates the files table in the database
func CreateFileTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS files (
		id SERIAL PRIMARY KEY,
		user_id INTEGER NOT NULL REFERENCES users(id),
		file_name VARCHAR(255) NOT NULL,
		upload_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		size BIGINT NOT NULL,
		s3_url VARCHAR(255) NOT NULL
	)`
	_, err := GetDB().Exec(query)
	return err
}