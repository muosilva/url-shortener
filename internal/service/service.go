package service

import (
	"context"
	"crypto/rand"
	"errors"
	"math/big"

	"github.com/muosilva/url-shortener/metrics"
)

const length = 8

var ErrURLNotFound = errors.New("url not found")

type Service interface {
	Create(ctx context.Context, url string) (string, error)
	Get(ctx context.Context, code string) (string, error)
}

type service struct {
	repo  URLRepository
	cache URLCache
}

func NewService(repo URLRepository, cache URLCache) Service {
	return &service{
		repo:  repo,
		cache: cache,
	}
}

func (s *service) Create(ctx context.Context, url string) (code string, err error) {
	code, err = generateCode(length)
	if err != nil {
		return "", err
	}

	err = s.repo.Create(ctx, code, url)
	if err != nil {
		return "", err
	}

	if err := s.cache.Set(ctx, code, url); err != nil {
		return "", err
	}

	metrics.URLsCreatedTotal.Inc()

	return code, nil
}

func (s *service) Get(ctx context.Context, code string) (longUrl string, err error) {
	result := metrics.ResultSuccess
	defer func() {
		metrics.URLLookupTotal.WithLabelValues(result).Inc()
	}()

	if len(code) == 0 || code == "" {
		err = errors.New("code is empty!")
		result = metrics.ResultError
		return "", err
	}

	longUrl, err = s.cache.Get(ctx, code)
	if err != nil {
		result = metrics.ResultError
		return "", err
	}
	if longUrl != "" {
		return longUrl, nil
	}

	longUrl, err = s.repo.GetByCode(ctx, code)
	if err != nil {
		result = metrics.ResultError
		return "", err
	}
	if longUrl == "" {
		result = metrics.ResultNotFound
		return "", ErrURLNotFound
	}

	_ = s.cache.Set(ctx, code, longUrl)

	return longUrl, nil
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
