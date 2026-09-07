package service

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
)

const length = 8

type Service interface {
	Create(ctx context.Context, url string) (string, error)
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
	fmt.Println(code)

	code, err = s.repo.Create(ctx, code, url)
	if err != nil {
		return "", err
	}

	if err := s.cache.Set(ctx, code, url); err != nil {
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
