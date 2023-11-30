package entity

import (
	"context"
)

type IURLService interface {
	CreateShortURL(ctx context.Context, url string) (string, error)
	GetLongURL(ctx context.Context, shortURL string) (string, error)
}

type IURLRepository interface {
	AddULR(ctx context.Context, key, value string) error
	GetLongURL(ctx context.Context, key string) (string, error)
}
