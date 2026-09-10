package service

import "context"

type URLRepository interface {
	Create(ctx context.Context, generatedCode, long_url string) error
	GetByCode(ctx context.Context, code string) (string, error)
}

type URLCache interface {
	Get(ctx context.Context, code string) (string, error)
	Set(ctx context.Context, code string, longURL string) error
}
