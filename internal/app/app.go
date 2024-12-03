package app

import (
	"fmt"
	"log/slog"
	"music_lib/internal/api"
	"music_lib/internal/api/handlers/song/add"
	"music_lib/internal/api/handlers/song/delete"
	"music_lib/internal/api/handlers/song/filter"
	"music_lib/internal/api/handlers/song/get"
	"music_lib/internal/api/handlers/song/update"
	"music_lib/internal/config"
	"music_lib/internal/dbs/postgres"
	"net/http"

	"github.com/swaggo/http-swagger"
	_ "music_lib/docs"

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

	err = postgres.Migrate(db)
	if err != nil {
		return err
	}
	slog.Debug("Migrations successfully applied")

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Route("/", func(r chi.Router) {
		song_svc := service.New(repo.NewRepository(db))
		client := api.NewClient(a.config.APIClient)

		r.Get("/get_song/{id}", get.New(song_svc))
		r.Get("/get_all_songs", get.NewAll(song_svc))
		r.Get("/filter_songs", filter.New(song_svc))
		r.Get("/get_song_text", get.NewText(song_svc))
		r.Post("/add", add.New(song_svc, client))
		r.Delete("/delete_song/{id}", delete.New(song_svc))
		r.Put("/update_song", update.New(song_svc))


		r.Get("/swagger/*", httpSwagger.WrapHandler)
	})

	addr := fmt.Sprintf("%s:%d", a.config.App.Host, a.config.App.Port)
	slog.Info("Starting server on", "address", fmt.Sprintf("http://%s", addr))
	srv := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	if err := srv.ListenAndServe(); err != nil {
		return err
	}

	return nil
}
