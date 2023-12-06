package rest

import (
	"fmt"
	controllerVideo "github.com/pranavraaghav/go-video-preview/src/delivery/rest/controllers/video"
	"github.com/pranavraaghav/go-video-preview/src/delivery/rest/routers"
	"github.com/pranavraaghav/go-video-preview/src/internal/usecase"
	"github.com/pranavraaghav/go-video-preview/src/utils"
	"net/http"
)

func StartNewRestDelivery(config *utils.Config, logger *utils.StandardLogger, usecases usecase.Usecases) {
	mux := http.NewServeMux()

	// Get all controllers
	videoController := controllerVideo.NewVideoController(logger, usecases.Video)

	// Set all routes
	routers.SetVideoRoutes(mux, videoController)

	serverPort := fmt.Sprintf(":%d", config.Port)

	// Start server
	logger.Infof("Starting server on port %s", serverPort)
	err := http.ListenAndServe(serverPort, mux)
	logger.Infof("Server stopped with error: %v", err)
}
