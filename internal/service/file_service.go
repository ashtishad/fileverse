package service

import (
	"context"
	"io"
	"log/slog"
	"strings"

	"github.com/ashtishad/fileverse/internal/domain"
	"github.com/ashtishad/fileverse/internal/infra/storage"
	"github.com/ashtishad/fileverse/pkg/utils"
)

// DefaultFileService is the primary implementation of FileService.
type DefaultFileService struct {
	repo   domain.FileRepository
	ipfs   *storage.IPFSStorage
	logger *slog.Logger
}

// NewFileService creates a new DefaultFileService with the given repository, IPFS storage, and logger.
func NewFileService(repo domain.FileRepository, ipfs *storage.IPFSStorage, logger *slog.Logger) *DefaultFileService {
	return &DefaultFileService{
		repo:   repo,
		ipfs:   ipfs,
		logger: logger,
	}
}

// SaveFile handles the logic for saving a new file.
func (s *DefaultFileService) SaveFile(ctx context.Context, fileName string, fileSize int64, fileReader io.Reader) (*domain.FileRespDTO, utils.APIError) { //nolint:lll
	// Upload to IPFS and get the hash
	cid, err := s.ipfs.UploadFile(fileReader)
	if err != nil {
		s.logger.Error("failed to upload file to IPFS", err)
		return nil, utils.InternalServerError("failed to upload file to IPFS", err)
	}

	file := domain.File{
		FileName: strings.TrimSpace(fileName),
		Size:     fileSize,
		IPFSHash: cid,
	}

	savedFile, apiErr := s.repo.SaveMeta(ctx, &file)
	if apiErr != nil {
		return nil, apiErr
	}

	fileRespDTO := domain.NewFileRespDTO(savedFile)

	return fileRespDTO, nil
}
