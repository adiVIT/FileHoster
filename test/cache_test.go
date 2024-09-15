package test

import (
	"testing"
	"filestore/internal/cache"
	"filestore/internal/files"
)

func TestSetAndGetFilesInCache(t *testing.T) {
	// Initialize Redis (ensure Redis is running and accessible)
	cache.InitRedis()

	// Sample file metadata
	userID := uint(1)
	fileMetadata := []files.FileMetadata{
		{ID: 1, UserID: userID, FileName: "test1.txt", Size: 1234, S3URL: "https://s3.amazonaws.com/bucket/test1.txt"},
		{ID: 2, UserID: userID, FileName: "test2.txt", Size: 5678, S3URL: "https://s3.amazonaws.com/bucket/test2.txt"},
	}

	// Set files in cache
	cache.SetFilesInCache(userID, fileMetadata)

	// Get files from cache
	cachedFiles, found := cache.GetFilesFromCache(userID)
	if !found {
		t.Fatalf("Expected files to be found in cache")
	}

	// Check if the cached files match the original
	if len(cachedFiles) != len(fileMetadata) {
		t.Errorf("Expected %d files, got %d", len(fileMetadata), len(cachedFiles))
	}

	for i, file := range cachedFiles {
		if file.ID != fileMetadata[i].ID || file.FileName != fileMetadata[i].FileName {
			t.Errorf("Expected file %v, got %v", fileMetadata[i], file)
		}
	}
}