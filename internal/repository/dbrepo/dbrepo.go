package dbrepo

import (
	"database/sql"
	"github.com/jjang65/booking-web-app/internal/config"
	"github.com/jjang65/booking-web-app/internal/repository"
)

type postgresDbRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

func NewPostgresRepo(conn *sql.DB, a *config.AppConfig) repository.DatabaseRepo {
	return &postgresDbRepo{
		App: a,
		DB:  conn,
	}
}
