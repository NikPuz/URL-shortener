package entity

import (
	"context"
)

type IURLService interface {
	CreateShortURL(ctx context.Context, url string) (string, error)
	GetURL(ctx context.Context, shortURL string) (string, error)
}

type IURLRepository interface {
	SetKeyValue(ctx context.Context, key, value string) error
	GetKeyValue(ctx context.Context, key string) (string, error)
}
