package main

import "fmt"

import (
	"music_lib/internal/app"
	"music_lib/internal/config"
	"music_lib/internal/logger"
)

func main() {
	cfg := config.MustLoad()

	logger.Setup()

	a := app.New(cfg)
	if err := a.Run(); err != nil {
		fmt.Println(err)
	}
}
