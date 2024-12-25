package models

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetadataMarshal(t *testing.T) {
	metadata := &Metadata{
		ID:        "1",
		UserID:    "user123",
		Name:      "file.txt",
		Path:      "/path/to/file.txt",
		MimeType:  "text/plain",
		URL:       "https://example.com/file.txt",
		CreatedAt: "2024-12-16T12:00:00Z",
		UpdatedAt: "2024-12-16T12:30:00Z",
	}

	jsonData, err := json.Marshal(metadata)
	assert.NoError(t, err)
	assert.NotNil(t, jsonData)

	expectedJSON := `{"id":"1","user_id":"user123","name":"file.txt","path":"/path/to/file.txt","mime_type":"text/plain","url":"https://example.com/file.txt","created_at":"2024-12-16T12:00:00Z","updated_at":"2024-12-16T12:30:00Z"}`
	assert.JSONEq(t, expectedJSON, string(jsonData))
}

func TestMetadataUnmarshal(t *testing.T) {
	jsonData := `{
		"id": "1",
		"user_id": "user123",
		"name": "file.txt",
		"path": "/path/to/file.txt",
		"mime_type": "text/plain",
		"url": "https://example.com/file.txt",
		"created_at": "2024-12-16T12:00:00Z",
		"updated_at": "2024-12-16T12:30:00Z"
	}`

	var metadata Metadata

	err := json.Unmarshal([]byte(jsonData), &metadata)
	assert.NoError(t, err)

	assert.Equal(t, "1", metadata.ID)
	assert.Equal(t, "user123", metadata.UserID)
	assert.Equal(t, "file.txt", metadata.Name)
	assert.Equal(t, "/path/to/file.txt", metadata.Path)
	assert.Equal(t, "text/plain", metadata.MimeType)
	assert.Equal(t, "https://example.com/file.txt", metadata.URL)
	assert.Equal(t, "2024-12-16T12:00:00Z", metadata.CreatedAt)
	assert.Equal(t, "2024-12-16T12:30:00Z", metadata.UpdatedAt)
}

func TestEmptyMetadata(t *testing.T) {
	var metadata Metadata

	jsonData, err := json.Marshal(metadata)
	assert.NoError(t, err)
	assert.Equal(t, `{"id":"","user_id":"","name":"","path":"","mime_type":"","url":"","created_at":"","updated_at":""}`, string(jsonData))

	err = json.Unmarshal([]byte(`{}`), &metadata)
	assert.NoError(t, err)
	assert.Equal(t, "", metadata.ID)
	assert.Equal(t, "", metadata.UserID)
	assert.Equal(t, "", metadata.Name)
	assert.Equal(t, "", metadata.Path)
	assert.Equal(t, "", metadata.MimeType)
	assert.Equal(t, "", metadata.URL)
	assert.Equal(t, "", metadata.CreatedAt)
	assert.Equal(t, "", metadata.UpdatedAt)
}
