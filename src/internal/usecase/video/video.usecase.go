package usecaseVideo

import (
	"archive/zip"
	"fmt"
	"github.com/pranavraaghav/go-video-preview/src/internal/domain/video"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

type VideoUsecaseImplementation struct{}

func NewVideoUsecaseImplementation() video.UseCase {
	return &VideoUsecaseImplementation{}
}

// GenerateImageZipFromVideo generates images from video
// returns path to zip file with all images
func (v *VideoUsecaseImplementation) GenerateImageZipFromVideo(filename string, width int, height int) (*string, error) {
	dirname := filename[:len(filename)-len(filepath.Ext(filename))] // filename without extension
	outputDirPath := fmt.Sprintf("./output/%s", dirname)
	err := os.MkdirAll(outputDirPath, os.ModePerm)
	if err != nil {
		return nil, err
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
		return nil, err
	}

	defer os.RemoveAll(outputDirPath)

	// Zip the directory
	err = os.MkdirAll("./zips", os.ModePerm)
	if err != nil {
		return nil, err
	}
	zipOutputFilePath := fmt.Sprintf("./zips/%s.zip", dirname)
	err = v.zipDirectory(outputDirPath, zipOutputFilePath)
	if err != nil {
		return nil, err
	}

	return &zipOutputFilePath, nil
}

func (v *VideoUsecaseImplementation) zipDirectory(source string, target string) error {
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
