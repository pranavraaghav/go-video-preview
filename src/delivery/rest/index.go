package rest

import (
	"fmt"
	controllerVideo "github.com/pranavraaghav/go-video-preview/src/delivery/rest/controllers/video"
	"github.com/pranavraaghav/go-video-preview/src/delivery/rest/routers"
	"github.com/pranavraaghav/go-video-preview/src/internal/usecase"
	"github.com/pranavraaghav/go-video-preview/src/utils"
	"net/http"
)

func StartNewRestDelivery(config *utils.Config, usecases usecase.Usecases) {
	mux := http.NewServeMux()

	// Get all controllers
	videoController := controllerVideo.NewVideoController(usecases.Video)

	// Set all routes
	routers.SetVideoRoutes(mux, videoController)

	serverPort := fmt.Sprintf(":%d", config.Port)

	// Start server
	http.ListenAndServe(serverPort, mux)
}
