package main

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleOk)
	mux.HandleFunc("/upload", handleUpload)
	mux.HandleFunc("/images", handleImages)

	const PORT = "8080"
	serverAddr := fmt.Sprintf(":%s", PORT)
	fmt.Printf("Listening on %s\n", serverAddr)
	http.ListenAndServe(serverAddr, mux)
}

func handleOk(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "OK\n")
}

type UploadResponse struct {
	Status   string `json:"status"`
	Message  string `json:"message"`
	Filename string `json:"filename"`
}

func handleUpload(w http.ResponseWriter, r *http.Request) {
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

func handleImages(w http.ResponseWriter, r *http.Request) {
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

	// dirname is name of file without extension
	dirname := filename[:len(filename)-len(filepath.Ext(filename))]
	outputDirPath := fmt.Sprintf("./output/%s", dirname)
	err := os.MkdirAll(outputDirPath, os.ModePerm)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Generate images from video using ffmpeg
	command := fmt.Sprintf(
		"~/ffmpeg -i ./uploads/%s -r 0.2 -s %dx%d -f image2 %s/%%03d.jpeg",
		filename,
		width,
		height,
		outputDirPath)

	err = exec.Command("sh", "-c", command).Run()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer os.RemoveAll(outputDirPath)

	// Zip the directory
	err = os.MkdirAll("./zips", os.ModePerm)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	zipOutputFileName := fmt.Sprintf("./zips/%s.zip", dirname)
	err = zipDirectory(outputDirPath, zipOutputFileName)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s.zip", dirname))
	w.WriteHeader(http.StatusOK)

	http.ServeFile(w, r, zipOutputFileName)
}

func zipDirectory(source string, target string) error {
	zipFile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	return filepath.Walk(source, func(filePath string, fileInfo os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories (they will be added automatically)
		if fileInfo.IsDir() {
			return nil
		}

		// Open the file to be zipped
		file, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer file.Close()

		// Create a new file header for the zip archive
		zipHeader, err := zip.FileInfoHeader(fileInfo)
		if err != nil {
			return err
		}

		// Specify the name of the file inside the zip archive
		zipHeader.Name, err = filepath.Rel(source, filePath)
		if err != nil {
			return err
		}

		// Create a file writer in the zip archive
		fileWriter, err := zipWriter.CreateHeader(zipHeader)
		if err != nil {
			return err
		}

		// Copy the file content to the zip archive
		_, err = io.Copy(fileWriter, file)
		return err
	})
}
