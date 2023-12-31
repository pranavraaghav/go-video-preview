package usecaseVideo

import (
	"archive/zip"
	"fmt"
	"github.com/pranavraaghav/go-video-preview/src/internal/domain/video"
	"github.com/pranavraaghav/go-video-preview/src/utils"
	"io"
	"math/rand"
	"mime/multipart"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

type videoUsecaseImplementation struct {
	ffmpegExecutablePath string
	logger               *utils.StandardLogger
}

func NewVideoUsecaseImplementation(config *utils.Config, logger *utils.StandardLogger) video.UseCase {
	ffmpegPath := config.FfmpegPath
	return &videoUsecaseImplementation{
		ffmpegExecutablePath: ffmpegPath,
		logger:               logger,
	}
}

func (v *videoUsecaseImplementation) GetMaxUploadFileSizeInBytes() int64 {
	return int64(4096 * 4096)
}

func (v *videoUsecaseImplementation) UploadFile(file multipart.File, fileHeader *multipart.FileHeader) (*string, error) {
	err := os.MkdirAll("uploads", os.ModePerm)
	if err != nil {
		return nil, err
	}

	fileExtension := filepath.Ext(fileHeader.Filename)
	dstFilename := fmt.Sprintf("%d%d", time.Now().UnixNano(), rand.Int63n(1000))
	dstFilePath := fmt.Sprintf("./uploads/%s%s", dstFilename, fileExtension)
	dst, err := os.Create(dstFilePath)
	if err != nil {
		return nil, err
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		return nil, err
	}

	return &dstFilePath, nil
}

// GenerateImageZipFromVideo generates images from video
// returns path to zip file with all images
func (v *videoUsecaseImplementation) GenerateImageZipFromVideo(
	filename string,
	width int,
	height int,
	intervalBetweenImages int,
) (*string, error) {
	dirname := filename[:len(filename)-len(filepath.Ext(filename))] // filename without extension
	outputDirPath := fmt.Sprintf("./output/%s", dirname)
	err := os.MkdirAll(outputDirPath, os.ModePerm)
	if err != nil {
		return nil, err
	}

	framesPerSecond := v.getFramesPerSecond(intervalBetweenImages)
	// Generate images from video using ffmpeg
	commandString := fmt.Sprintf(
		"%s -i ./uploads/%s -r %0.2f -s %dx%d -f image2 %s/%%03d.jpeg",
		v.ffmpegExecutablePath,
		filename,
		framesPerSecond,
		width,
		height,
		outputDirPath,
	)
	err = exec.Command("sh", "-c", commandString).Run()
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

func (v *videoUsecaseImplementation) zipDirectory(source string, target string) error {
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

func (v *videoUsecaseImplementation) getFramesPerSecond(interval int) float32 {
	return 1 / float32(interval)
}
