package postgres

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"music_lib/internal/config"
)

func New(cfg config.PostgresConfig) (*sqlx.DB, error) {
	dataSource := fmt.Sprintf("user=%s password=%s host=%s port=%d dbname=%s sslmode=disable", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name)

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	timeoutExceed := time.After(cfg.ConnTimeExceed)
	for {
		select {
		case <-timeoutExceed:
			return nil, errors.New("Connection time exceed")
		case <-ticker.C:

			if conn, err := sqlx.Connect("postgres", dataSource); err == nil {
				if err = conn.Ping(); err == nil {
					return conn, err
				}
			}
		}
	}
}
