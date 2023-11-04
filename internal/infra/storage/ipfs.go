package storage

import (
	"errors"
	"fmt"
	"io"

	shell "github.com/ipfs/go-ipfs-api"
)

type IPFSStorage struct {
	shell *shell.Shell
}

func NewIPFSStorage(ipfsAddress string) *IPFSStorage {
	return &IPFSStorage{
		shell: shell.NewShell(ipfsAddress),
	}
}

func (s *IPFSStorage) UploadFile(r io.Reader) (string, error) {
	cid, err := s.shell.Add(r)
	if err != nil {
		return "", fmt.Errorf("error while adding file to ipfs: %w", err)
	}

	return cid, nil
}

// RetrieveFileChunk reads a specific chunk from an IPFS-stored file given its CID, offset, and size.
// It returns the requested file segment as a byte slice or an error if the read is unsuccessful.
func (s *IPFSStorage) RetrieveFileChunk(cid string, offset int64, size int64) ([]byte, error) {
	reader, err := s.shell.Cat(cid)
	if err != nil {
		return nil, fmt.Errorf("error while retrieving file from IPFS: %w", err)
	}
	defer reader.Close()

	if offset > 0 {
		_, err = io.CopyN(io.Discard, reader, offset)
		if err != nil {
			return nil, fmt.Errorf("error while skipping to offset in file: %w", err)
		}
	}

	buf := make([]byte, size)
	n, err := io.ReadFull(reader, buf)

	if err != nil && err != io.EOF && !errors.Is(err, io.ErrUnexpectedEOF) {
		return nil, fmt.Errorf("error while reading chunk from file: %w", err)
	}

	return buf[:n], nil
}
