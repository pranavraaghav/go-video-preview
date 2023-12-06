package rest

import (
	controllerVideo "github.com/pranavraaghav/go-video-preview/src/delivery/rest/controllers/video"
	"github.com/pranavraaghav/go-video-preview/src/delivery/rest/routers"
	"github.com/pranavraaghav/go-video-preview/src/internal/usecase"
	"net/http"
)

func StartNewRestDelivery(usecases usecase.Usecases) {
	mux := http.NewServeMux()

	// Get all controllers
	videoController := controllerVideo.NewVideoController(usecases.Video)

	// Set all routes
	routers.SetVideoRoutes(mux, videoController)

	serverPort := ":8080"

	// Start server
	err := http.ListenAndServe(serverPort, mux)
	if err != nil {
		panic(err)
	}
}
