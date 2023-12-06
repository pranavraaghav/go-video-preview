package routers

import (
	controllerVideo "github.com/pranavraaghav/go-video-preview/src/delivery/rest/controllers/video"
	"net/http"
)

func SetVideoRoutes(mux *http.ServeMux, c controllerVideo.VideoController) {
	mux.HandleFunc("/upload", c.HandleUpload)
	mux.HandleFunc("/images", c.HandleImages)
}
