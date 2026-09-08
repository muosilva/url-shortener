package postgres

import (
	"context"
	"database/sql"
	"errors"
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

func (r *PostgresURLRepository) GetByCode(ctx context.Context, code string) (longUrl string, err error) {
	err = r.db.QueryRowContext(ctx, `SELECT long_url FROM urls WHERE code = $1`, code).Scan(&longUrl)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil
		}

		return "", err
	}

	return longUrl, nil
}
