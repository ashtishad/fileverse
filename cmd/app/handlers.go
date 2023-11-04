package app

import (
	"net/http"

	"github.com/ashtishad/fileverse/internal/service"
	"github.com/gin-gonic/gin"
)

type FileHandlers struct {
	s *service.DefaultFileService
}

// SaveFileHandler handles the HTTP request for saving file metadata.
func (fh *FileHandlers) SaveFileHandler(c *gin.Context) {
	if err := c.Request.ParseMultipartForm(32 << 20); err != nil { // 32 MB max memory
		c.JSON(http.StatusBadRequest, gin.H{"error": "File upload error: " + err.Error()})
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File retrieval error: " + err.Error()})
		return
	}
	defer file.Close()

	fileName := header.Filename
	fileSize := header.Size

	// Pass the file along with its metadata to the service layer
	fileResp, apiErr := fh.s.SaveFile(c.Request.Context(), fileName, fileSize, file)
	if apiErr != nil {
		c.JSON(apiErr.Code(), gin.H{
			"error": apiErr.Error(),
		})

		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"file": fileResp,
	})
}
