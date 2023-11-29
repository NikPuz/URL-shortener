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

func (r URLRepository) SetKeyValue(ctx context.Context, key, value string) error {

	client := gossdb.Client()
	defer client.Close()

	err := client.Set(key, value)
	if err != nil {
		return err
	}

	return nil
}

func (r URLRepository) GetKeyValue(ctx context.Context, key string) (string, error) {

	client := gossdb.Client()
	defer client.Close()

	value, err := client.Get(key)
	if err != nil {
		return "", err
	}

	return value.String(), nil
}
