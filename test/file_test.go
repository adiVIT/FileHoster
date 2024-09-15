package test

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
	"filestore/internal/handlers"
	"filestore/internal/auth"
)

func TestUploadFileHandler(t *testing.T) {
	// Create a new file upload request
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "test.txt")
	if err != nil {
		t.Fatalf("Failed to create form file: %v", err)
	}
	part.Write([]byte("This is a test file"))
	writer.Close()

	// Create a new HTTP request
	req, err := http.NewRequest("POST", "/upload", body)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Mock authentication
	token, _ := auth.GenerateJWT(1)
	req.Header.Set("Authorization", token)

	// Record the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.UploadFileHandler)
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body
	expected := `{"url":"`
	if !bytes.Contains(rr.Body.Bytes(), []byte(expected)) {
		t.Errorf("Handler returned unexpected body: got %v", rr.Body.String())
	}
}

func TestListFilesHandler(t *testing.T) {
	// Create a new HTTP request
	req, err := http.NewRequest("GET", "/files", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Mock authentication
	token, _ := auth.GenerateJWT(1)
	req.Header.Set("Authorization", token)

	// Record the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.ListFilesHandler)
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body
	expected := `[]` // Assuming no files initially
	if rr.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}