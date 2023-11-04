package domain

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/ashtishad/fileverse/pkg/utils"
)

type FileRepoDB struct {
	db *sql.DB
	l  *slog.Logger
}

func NewFileRepoDB(db *sql.DB, l *slog.Logger) *FileRepoDB {
	return &FileRepoDB{
		db: db,
		l:  l,
	}
}

// checkExists checks if a file with the given IPFS hash already exists in the database.
func (d *FileRepoDB) checkExists(ctx context.Context, ipfsHash string) bool {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM file_metadata WHERE ipfs_hash=$1)`
	err := d.db.QueryRowContext(ctx, query, ipfsHash).Scan(&exists)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		d.l.Error("error checking for existing IPFS hash", "err", err)
	}

	return exists
}

// SaveMeta adds a new file's metadata to the database after checking for existing IPFS hash.
func (d *FileRepoDB) SaveMeta(ctx context.Context, file *File) (*File, utils.APIError) {
	if d.checkExists(ctx, file.IPFSHash) {
		return nil, utils.ConflictError(fmt.Sprintf("a file with the IPFS hash '%s' already exists", file.IPFSHash))
	}

	const sqlInsertFile = `INSERT INTO file_metadata (file_name, size, ipfs_hash) 
						   VALUES ($1, $2, $3) RETURNING id, file_id, timestamp;`

	var newID int
	var newFileID string
	var newTimestamp time.Time

	err := d.db.QueryRowContext(ctx, sqlInsertFile,
		file.FileName, file.Size, file.IPFSHash).Scan(&newID, &newFileID, &newTimestamp)
	if err != nil {
		d.l.Error("error inserting new file record", "err", err)

		if errors.Is(err, sql.ErrNoRows) {
			return nil, utils.NotFoundError("no new file record was created")
		}

		return nil, utils.InternalServerError("error inserting new file record", err)
	}

	file.ID = newID
	file.FileID = newFileID
	file.Timestamp = newTimestamp

	return file, nil
}
