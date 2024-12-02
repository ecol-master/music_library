package app

import (
	"fmt"
	"log/slog"
	"music_lib/internal/api"
	"music_lib/internal/api/handlers/song/add"
	"music_lib/internal/api/handlers/song/delete"
	"music_lib/internal/config"
	"music_lib/internal/dbs/postgres"
	"net/http"

	repo "music_lib/internal/repositories/song"
	service "music_lib/internal/services/song"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type App struct {
	config *config.Config
}

func New(config *config.Config) *App {
	return &App{
		config: config,
	}
}

func (a *App) Run() error {
	// Connectin to database
	db, err := postgres.New(a.config.Postgres)
	if err != nil {
		return err
	}
	slog.Debug("Suceessfully connected to database")

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Route("/", func(r chi.Router) {
		song_svc := service.New(repo.NewRepository(db))
		client := api.NewClient(a.config.APIClient)

		r.Post("/add", add.New(song_svc, client))
		r.Post("/delete", delete.New(song_svc))
	})

	addr := fmt.Sprintf("%s:%d", a.config.App.Host, a.config.App.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	if err := srv.ListenAndServe(); err != nil {
		return err
	}

	return nil
}
