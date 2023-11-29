package service

import (
	"context"
	"crypto/sha256"
	"github.com/jxskiss/base62"
	"url-shortner/internal/entity"
)

type URLService struct {
	URLRepo entity.IURLRepository
}

func NewURLService(urlRepo entity.IURLRepository) entity.IURLService {
	URLService := new(URLService)
	URLService.URLRepo = urlRepo
	return URLService
}

func (s URLService) CreateShortURL(ctx context.Context, url string) (string, error) {
	sha := sha256.New()
	sha.Write([]byte(url))
	hash := sha.Sum(nil)

	shortURL := string(base62.Encode(hash)[:8])

	err := s.URLRepo.SetKeyValue(ctx, shortURL, url)
	if err != nil {
		return "", err
	}

	return shortURL, nil
}

func (s URLService) GetURL(ctx context.Context, shortURL string) (string, error) {

	url, err := s.URLRepo.GetKeyValue(ctx, shortURL)
	if err != nil {
		return "", err
	}

	return url, nil
}
