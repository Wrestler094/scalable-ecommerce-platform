package postgres

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Postgres struct {
	DB *sqlx.DB
}

func NewConnect(url string) (*Postgres, error) {
	db, err := sqlx.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &Postgres{DB: db}, nil
}

func (p *Postgres) Close() error {
	return p.DB.Close()
}
