package main

import (
	"github.com/pranavraaghav/go-video-preview/src/delivery/rest"
	"github.com/pranavraaghav/go-video-preview/src/internal/usecase"
)

func main() {
	usecases := usecase.InitUsecases()
	rest.StartNewRestDelivery(usecases)
}
