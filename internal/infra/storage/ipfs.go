package storage

import (
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

func (s *IPFSStorage) RetrieveFile(cid string) (io.ReadCloser, error) {
	reader, err := s.shell.Cat(cid)
	if err != nil {
		return nil, fmt.Errorf("error while retriving file from ipfs: %w", err)
	}

	return reader, nil
}
