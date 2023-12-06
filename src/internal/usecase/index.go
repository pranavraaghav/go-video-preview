package usecase

import (
	"github.com/pranavraaghav/go-video-preview/src/internal/domain/video"
	usecaseVideo "github.com/pranavraaghav/go-video-preview/src/internal/usecase/video"
	"github.com/pranavraaghav/go-video-preview/src/utils"
)

type Usecases struct {
	Video video.UseCase
}

func InitUsecases(config *utils.Config, logger *utils.StandardLogger) Usecases {
	videoUsecase := usecaseVideo.NewVideoUsecaseImplementation(config, logger)
	return Usecases{
		Video: videoUsecase,
	}
}
