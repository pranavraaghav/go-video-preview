package main

import (
	"github.com/pranavraaghav/go-video-preview/src/delivery/rest"
	"github.com/pranavraaghav/go-video-preview/src/internal/usecase"
	"github.com/pranavraaghav/go-video-preview/src/utils"
)

func main() {
	config := utils.NewConfig()
	logger := utils.NewLogger(config)

	usecases := usecase.InitUsecases(config, logger)
	rest.StartNewRestDelivery(config, logger, usecases)
}
