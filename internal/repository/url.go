package repository

import (
	"context"
	"github.com/seefan/gossdb/v2"
	"url-shortner/internal/entity"
)

type URLRepository struct {
}

func NewURLRepository() entity.IURLRepository {
	URLRepository := new(URLRepository)
	return URLRepository
}

func (r URLRepository) AddULR(ctx context.Context, shortURL, LongURL string) error {

	client := gossdb.Client()
	defer client.Close()

	err := client.HSet("Links", shortURL, LongURL)
	if err != nil {
		return err
	}

	return nil
}

func (r URLRepository) GetLongURL(ctx context.Context, shortURL string) (string, error) {

	client := gossdb.Client()
	defer client.Close()

	value, err := client.HGet("Links", shortURL)
	if err != nil {
		return "", err
	}

	return value.String(), nil
}
