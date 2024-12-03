package main

import "fmt"

import (
	_ "music_lib/docs"
	"music_lib/internal/app"
	"music_lib/internal/config"
	"music_lib/internal/logger"
)

// @title Music Library API
// @version 1.0
// @host localhost:8080
// @BasePath /

func main() {
	cfg := config.MustLoad()

	logger.Setup()

	a := app.New(cfg)
	if err := a.Run(); err != nil {
		fmt.Println(err)
	}
}
