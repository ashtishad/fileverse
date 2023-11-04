package domain

import (
	"context"

	"github.com/ashtishad/fileverse/pkg/utils"
)

type FileRepository interface {
	SaveMeta(ctx context.Context, file *File) (*File, utils.APIError)
	FindMeta(ctx context.Context, fileID string) (*File, utils.APIError)
}
