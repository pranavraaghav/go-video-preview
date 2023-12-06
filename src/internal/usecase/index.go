package usecase

import (
	"github.com/pranavraaghav/go-video-preview/src/internal/domain/video"
	usecaseVideo "github.com/pranavraaghav/go-video-preview/src/internal/usecase/video"
	"github.com/pranavraaghav/go-video-preview/src/utils"
)

type Usecases struct {
	Video video.UseCase
}

func InitUsecases(config *utils.Config) Usecases {
	videoUsecase := usecaseVideo.NewVideoUsecaseImplementation(config.FfmpegPath)
	return Usecases{
		Video: videoUsecase,
	}
}
