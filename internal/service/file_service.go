package service

import (
	"context"
	"log/slog"
	"strings"

	"github.com/ashtishad/fileverse/internal/domain"
	"github.com/ashtishad/fileverse/pkg/utils"
)

// FileService defines the interface for file service operations.
type FileService interface {
	SaveFile(ctx context.Context, fileData domain.NewFileReqDTO) (*domain.FileRespDTO, utils.APIError)
}

// DefaultFileService is the default implementation of FileService.
type DefaultFileService struct {
	repo domain.FileRepository
	l    *slog.Logger
}

// NewFileService creates a new DefaultFileService with the given repository and logger.
func NewFileService(repo domain.FileRepository, l *slog.Logger) *DefaultFileService {
	return &DefaultFileService{
		repo: repo,
		l:    l,
	}
}

// SaveFile handles the logic for saving a new file.
func (s *DefaultFileService) SaveFile(ctx context.Context, fileData domain.NewFileReqDTO) (*domain.FileRespDTO, utils.APIError) { //nolint:lll
	// ToDo: Validate file input data
	file := domain.File{
		FileName: strings.TrimSpace(fileData.FileName),
		Size:     1000,         // dummy for now, will be generated from file
		IPFSHash: "samplehash", // dummy for now, will be generated from file
	}

	savedFile, apiErr := s.repo.SaveMeta(ctx, &file)
	if apiErr != nil {
		return nil, apiErr
	}

	fileRespDTO := domain.NewFileRespDTO(savedFile)

	return fileRespDTO, nil
}
