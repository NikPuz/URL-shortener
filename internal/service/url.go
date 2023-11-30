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
	var shortURL string
	hashObject := url

loop:
	for {
		sha := sha256.New()
		sha.Write([]byte(hashObject))
		hash := sha.Sum(nil)

		shortURL = string(base62.Encode(hash)[:8])

		longURL, err := s.URLRepo.GetLongURL(ctx, shortURL)
		if err != nil {
			return "", err
		}

		switch longURL {
		case url: // Нашли, возвращаем
			return shortURL, nil
		case "": // Не нашли, записываем
			break loop
		default: // Коллизия, перехэшиваем
			hashObject = shortURL
		}
	}

	err := s.URLRepo.AddULR(ctx, shortURL, url)
	if err != nil {
		return "", err
	}

	return shortURL, nil
}

func (s URLService) GetLongURL(ctx context.Context, shortURL string) (string, error) {
	return s.URLRepo.GetLongURL(ctx, shortURL)
}
