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

// RetrieveFile handles the logic for retrieving an existing file.
func (s *DefaultFileService) RetrieveFile(ctx context.Context, fileID string) ([]byte, utils.APIError) {
	fileMeta, apiErr := s.repo.FindMeta(ctx, fileID)
	if apiErr != nil {
		s.logger.Error("failed to get file metadata", "fileID", fileID, "error", apiErr)
		return nil, apiErr
	}

	fileReader, err := s.ipfs.RetrieveFile(fileMeta.IPFSHash)
	if err != nil {
		s.logger.Error("failed to retrieve file from IPFS", "hash", fileMeta.IPFSHash, "error", err)
		return nil, utils.InternalServerError("failed to retrieve file from IPFS", err)
	}
	defer fileReader.Close()

	fileContent, err := io.ReadAll(fileReader)
	if err != nil {
		s.logger.Error("failed to read file content", "hash", fileMeta.IPFSHash, "error", err)
		return nil, utils.InternalServerError("failed to read file content", err)
	}

	return fileContent, nil
}
