package db

import (
	"adt/config"
	"log/slog"

	"github.com/jmoiron/sqlx"
)

func NewMariaDB(cfg config.Config) (*sqlx.DB, error) {
	slog.Info("Connecting to MariaDB . . . ")
	db, err := sqlx.Connect(cfg.DBDriver, cfg.DBSource)
	if err != nil {
		return nil, err
	}

	return db, nil
}
