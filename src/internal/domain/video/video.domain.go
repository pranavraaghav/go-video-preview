package video

import "mime/multipart"

type UseCase interface {
	GenerateImageZipFromVideo(filename string, width int, height int) (*string, error)
	UploadFile(file multipart.File, fileHeader *multipart.FileHeader) (*string, error)
	GetMaxUploadFileSizeInBytes() int64
}
