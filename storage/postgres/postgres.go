package postgres

import (
	"context"
	"crypto/rand"
	"database/sql"
	"math/big"
)

const codeLength = 8

type PostgresURLRepository struct {
	db *sql.DB
}

func NewPostgresURLRepository(db *sql.DB) *PostgresURLRepository {
	return &PostgresURLRepository{
		db: db,
	}
}

func (r *PostgresURLRepository) Create(ctx context.Context, url string) (string, error) {
	code, err := generateCode(codeLength)
	if err != nil {
		return "", err
	}

	_, err = r.db.ExecContext(
		ctx,
		`INSERT INTO urls (code, long_url) VALUES ($1, $2)`,
		code,
		url,
	)
	if err != nil {
		return "", err
	}

	return code, nil
}

func generateCode(length int) (string, error) {
	const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	code := make([]byte, length)
	for i := range code {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(alphabet))))
		if err != nil {
			return "", err
		}

		code[i] = alphabet[n.Int64()]
	}

	return string(code), nil
}
