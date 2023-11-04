package app

import (
	"net/http"

	"github.com/ashtishad/fileverse/internal/domain"
	"github.com/ashtishad/fileverse/internal/service"
	"github.com/gin-gonic/gin"
)

type FileHandlers struct {
	s service.FileService
}

// SaveFileHandler handles the HTTP request for saving file metadata.
func (fh *FileHandlers) SaveFileHandler(c *gin.Context) {
	var fileDataRequest domain.NewFileReqDTO
	if err := c.ShouldBindJSON(&fileDataRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fileResp, apiErr := fh.s.SaveFile(c.Request.Context(), fileDataRequest)
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
