package domain

import "time"

// File represents the file metadata that will be stored in the postgres database.
type File struct {
	ID        int       `json:"id"`        // Corresponds to SERIAL PRIMARY KEY
	FileID    string    `json:"fileId"`    // UUID from the database, auto-generated
	FileName  string    `json:"fileName"`  // Name of the file
	Size      int64     `json:"size"`      // File size in bytes
	Timestamp time.Time `json:"timestamp"` // Timestamp of when the file was added, auto-generated
	IPFSHash  string    `json:"ipfsHash"`  // The IPFS hash corresponding to the file
}
