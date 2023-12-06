package utils

import (
	"fmt"
	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
)

type Config struct {
	FfmpegPath string `env:"PATH_TO_FFMPEG_EXECUTABLE" envDefault:""`
	Port       int    `env:"PORT" envDefault:"3000"`
}

func NewConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		return nil
	}
	cfg := Config{}
	if err = env.Parse(&cfg); err != nil {
		fmt.Println("Error parsing env variables")
		panic(err)
	}
	return &cfg
}
