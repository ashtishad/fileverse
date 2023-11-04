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

// GetFileHandler handles the HTTP request for retrieving file content.
func (fh *FileHandlers) GetFileHandler(c *gin.Context) {
	fileID := c.Param("fileId")

	fileContent, apiErr := fh.s.RetrieveFile(c.Request.Context(), fileID)
	if apiErr != nil {
		c.JSON(apiErr.Code(), gin.H{
			"error": apiErr.Error(),
		})

		return
	}

	c.Writer.WriteHeader(http.StatusOK)

	if _, err := c.Writer.Write(fileContent); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to write file content to response: " + err.Error(),
		})
	}
}
