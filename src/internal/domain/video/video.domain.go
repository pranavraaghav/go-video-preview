package video

type UseCase interface {
	GenerateImageZipFromVideo(filename string, width int, height int) (*string, error)
}
