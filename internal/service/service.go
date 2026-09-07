package service

import (
	"context"
	"encoding/base32"
	"fmt"
)

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
	hashed := hashUrl(url)
	fmt.Println("Hashed url: ", hashed)

	code, err = s.repo.Create(ctx, hashed)
	if err != nil {
		return "", err
	}

	if err := s.cache.Set(ctx, code, url); err != nil {
		return "", err
	}

	return code, nil
}

func hashUrl(url string) (hashed string) {
	hashed = base32.StdEncoding.EncodeToString([]byte(url))
	return hashed
}
