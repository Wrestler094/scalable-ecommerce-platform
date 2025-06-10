package postgres

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Postgres struct {
	DB *sqlx.DB
}

func NewConnect(url string) (*Postgres, error) {
	const op = "postgres.NewConnect"

	db, err := sqlx.Open("postgres", url)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to open connection: %w", op, err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("%s: failed to ping database: %w", op, err)
	}

	return &Postgres{DB: db}, nil
}

func (p *Postgres) Close() error {
	const op = "postgres.Close"

	if err := p.DB.Close(); err != nil {
		return fmt.Errorf("%s: failed to close db: %w", op, err)
	}

	return nil
}
