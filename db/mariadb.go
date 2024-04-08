package db

import (
	"adt/config"

	"github.com/jmoiron/sqlx"
)

func NewMariaDB(cfg config.Config) (*sqlx.DB, error) {
	db, err := sqlx.Connect(cfg.DBDriver, cfg.DBSource)
	if err != nil {
		return nil, err
	}

	return db, nil
}
