package app

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/ashtishad/fileverse/internal/domain"
	"github.com/ashtishad/fileverse/internal/infra/database"
	"github.com/ashtishad/fileverse/internal/infra/storage"
	"github.com/ashtishad/fileverse/internal/service"
	"github.com/ashtishad/fileverse/pkg/utils"
	"github.com/gin-gonic/gin"
)

func Start() {
	gin.SetMode(gin.ReleaseMode)

	logger := utils.InitSlogger()

	utils.SanityCheck(logger)

	dbClient := database.GetDBClient(logger)

	ipfsClient := storage.NewIPFSStorage("localhost:5001")

	fileRepositoryDB := domain.NewFileRepoDB(dbClient, logger)

	fileService := service.NewFileService(fileRepositoryDB, ipfsClient, logger)

	router := gin.Default()

	fileHandlers := FileHandlers{s: fileService}

	router.POST("/upload", fileHandlers.SaveFileHandler)

	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", os.Getenv("API_HOST"), os.Getenv("API_PORT")),
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	logger.Info("Starting server " + srv.Addr)

	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Error("could not start server", "err", err)
	}
}
