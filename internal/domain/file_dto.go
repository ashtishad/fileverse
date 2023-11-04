package domain

import "time"

// FileRespDTO represents the response data to be sent back to the client when they query for file metadata.
type FileRespDTO struct {
	FileID    string    `json:"fileId"`    // UUID as a string.
	FileName  string    `json:"fileName"`  // Original name of the file.
	Size      int64     `json:"size"`      // Size of the file in bytes.
	Timestamp time.Time `json:"timestamp"` // Time when the file metadata was created.
}

// NewFileRespDTO creates a new DTO from a File domain model.
func NewFileRespDTO(file *File) *FileRespDTO {
	return &FileRespDTO{
		FileID:    file.FileID,
		FileName:  file.FileName,
		Size:      file.Size,
		Timestamp: file.Timestamp,
	}
}
