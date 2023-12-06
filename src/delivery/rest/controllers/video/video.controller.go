package controllerVideo

import (
	"encoding/json"
	"fmt"
	"github.com/pranavraaghav/go-video-preview/src/internal/domain/video"
	"net/http"
	"path/filepath"
	"strconv"
)

type VideoController interface {
	HandleUpload(w http.ResponseWriter, r *http.Request)
	HandleImages(w http.ResponseWriter, r *http.Request)
}

type videoControllerImplementation struct {
	usecase video.UseCase
}

func NewVideoController(usecase video.UseCase) VideoController {
	return &videoControllerImplementation{
		usecase: usecase,
	}
}

func (v *videoControllerImplementation) HandleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	maxUploadFileSizeInBytes := v.usecase.GetMaxUploadFileSizeInBytes()
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadFileSizeInBytes)
	if err := r.ParseMultipartForm(maxUploadFileSizeInBytes); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	storedFilePath, err := v.usecase.UploadFile(file, fileHeader)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	storedFileName := filepath.Base(*storedFilePath)
	p := UploadResponsePayload{
		Filename: storedFileName,
	}
	err = json.NewEncoder(w).Encode(p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (v *videoControllerImplementation) HandleImages(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get filename and dimensions from query params
	filename := r.URL.Query().Get("filename")
	if filename == "" {
		http.Error(w, "filename not provided in query params", http.StatusBadRequest)
		return
	}
	height, err1 := strconv.Atoi(r.URL.Query().Get("height"))
	width, err2 := strconv.Atoi(r.URL.Query().Get("width"))
	if err1 != nil || err2 != nil {
		http.Error(w, "Invalid height or width", http.StatusBadRequest)
		return
	}

	// Generate zip with images from video
	zipOutputFilePath, err := v.usecase.GenerateImageZipFromVideo(filename, width, height)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Serve zip file
	zipOutputFileName := filepath.Base(*zipOutputFilePath)
	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", zipOutputFileName))
	w.WriteHeader(http.StatusOK)
	http.ServeFile(w, r, *zipOutputFilePath)
}
