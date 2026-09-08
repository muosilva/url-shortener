package postgres

import (
	"context"
	"database/sql"
)

type PostgresURLRepository struct {
	db *sql.DB
}

func NewPostgresURLRepository(db *sql.DB) *PostgresURLRepository {
	return &PostgresURLRepository{
		db: db,
	}
}

func (r *PostgresURLRepository) Create(ctx context.Context, generateCode, long_url string) (err error) {
	_, err = r.db.ExecContext(
		ctx,
		`INSERT INTO urls (code, long_url) VALUES ($1, $2)`,
		generateCode,
		long_url,
	)
	if err != nil {
		return err
	}

	return nil
}
