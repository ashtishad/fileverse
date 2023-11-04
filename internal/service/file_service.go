package service

import (
	"context"
	"io"
	"log/slog"
	"math"
	"strings"
	"sync"

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

// RetrieveFile fetches the contents of a file identified by fileID in concurrent file chunks,
// Leverages goroutines to retrieve each chunk from IPFS and assembles them into a single byte slice in correct order,
// If successful, it returns the complete file content as []byte, otherwise an error.
func (s *DefaultFileService) RetrieveFile(ctx context.Context, fileID string) ([]byte, utils.APIError) {
	fileMeta, apiErr := s.repo.FindMeta(ctx, fileID)
	if apiErr != nil {
		s.logger.Error("failed to get file metadata", "fileID", fileID, "error", apiErr)
		return nil, apiErr
	}

	numChunks := int(math.Ceil(float64(fileMeta.Size) / float64(utils.FileChunkSize)))

	fileContent := make([]byte, fileMeta.Size)
	var wg sync.WaitGroup
	wg.Add(numChunks)

	for i := 0; i < numChunks; i++ {
		go func(chunkNum int) {
			defer wg.Done()

			offset := chunkNum * utils.FileChunkSize
			size := utils.FileChunkSize

			if chunkNum == numChunks-1 {
				size = int(fileMeta.Size) - offset
			}

			chunk, err := s.ipfs.RetrieveFileChunk(fileMeta.IPFSHash, int64(offset), int64(size))
			if err != nil {
				s.logger.Error("failed to retrieve file chunk from IPFS", "chunkNum", chunkNum, "error", err)
				return
			}

			copy(fileContent[offset:], chunk)
		}(i)
	}

	wg.Wait()

	return fileContent, nil
}
