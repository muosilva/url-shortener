package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/muosilva/url-shortener/metrics"
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
	start := time.Now()
	result := metrics.ResultSuccess
	defer func() {
		metrics.ObserveDBOperation("insert_url", result, time.Since(start))
	}()

	_, err = r.db.ExecContext(
		ctx,
		`INSERT INTO urls (code, long_url) VALUES ($1, $2)`,
		generateCode,
		long_url,
	)
	if err != nil {
		result = metrics.ResultError
		return err
	}

	return nil
}

func (r *PostgresURLRepository) GetByCode(ctx context.Context, code string) (longUrl string, err error) {
	start := time.Now()
	result := metrics.ResultSuccess
	defer func() {
		metrics.ObserveDBOperation("find_by_code", result, time.Since(start))
	}()

	err = r.db.QueryRowContext(ctx, `SELECT long_url FROM urls WHERE code = $1`, code).Scan(&longUrl)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			result = metrics.ResultNotFound
			return "", nil
		}

		result = metrics.ResultError
		return "", err
	}

	return longUrl, nil
}
