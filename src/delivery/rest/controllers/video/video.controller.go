package controllerVideo

import (
	"encoding/json"
	"fmt"
	"github.com/pranavraaghav/go-video-preview/src/internal/domain/video"
	"github.com/pranavraaghav/go-video-preview/src/internal/usecase"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
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
	MAX_UPLOAD_SIZE := int64(4096 * 4096)
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, MAX_UPLOAD_SIZE)
	if err := r.ParseMultipartForm(MAX_UPLOAD_SIZE); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	defer file.Close()

	err = os.MkdirAll("uploads", os.ModePerm)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	fileExtension := filepath.Ext(fileHeader.Filename)
	dstFilename := fmt.Sprintf("%d%d", time.Now().UnixNano(), rand.Int63n(1000))
	dst, err := os.Create(fmt.Sprintf("./uploads/%s%s", dstFilename, fileExtension))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	p := UploadResponse{
		Status:   "success",
		Message:  "File uploaded successfully",
		Filename: fmt.Sprintf("%s%s", dstFilename, fileExtension),
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

	// Query params
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

	usecases := usecase.InitUsecases()
	zipOutputFilePath, err := usecases.Video.GenerateImageZipFromVideo(filename, width, height)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	zipOutputFileName := filepath.Base(*zipOutputFilePath)
	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", zipOutputFileName))
	w.WriteHeader(http.StatusOK)

	http.ServeFile(w, r, *zipOutputFilePath)
}
