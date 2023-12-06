package utils

import (
	"fmt"
	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
)

type Config struct {
	FfmpegPath string `env:"PATH_TO_FFMPEG_EXECUTABLE" envDefault:""`
}

func NewConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		return nil
	}
	cfg := Config{}
	if err = env.Parse(&cfg); err != nil {
		fmt.Println("Error parsing .env file")
		panic(err)
	}
	return &cfg
}
