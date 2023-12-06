package usecase

import (
	"github.com/pranavraaghav/go-video-preview/src/internal/domain/video"
	usecaseVideo "github.com/pranavraaghav/go-video-preview/src/internal/usecase/video"
)

type Usecases struct {
	Video video.UseCase
}

func InitUsecases() Usecases {
	videoUsecase := usecaseVideo.NewVideoUsecaseImplementation()
	return Usecases{
		Video: videoUsecase,
	}
}
