package service

import (
	"context"
	"time"

	"github.com/NghiaLeopard/bookmark-management/internal/repository"
)

type ShortenUrlService interface {
	CreateShortenUrl(ctx context.Context, url string, expire time.Duration) (string, error)
	GetUrlByCode(ctx context.Context, code string) (string, error)
}

type shortenUrlService struct {
	urlStorage repository.UrlStorage
	genCode    GenPass
}

func NewShortenUrlService(urlStorage repository.UrlStorage, genCode GenPass) ShortenUrlService {
	return &shortenUrlService{
		urlStorage: urlStorage,
		genCode:    genCode,
	}
}

var length = 7

func (s *shortenUrlService) CreateShortenUrl(ctx context.Context, url string, expire time.Duration) (string, error) {
	code, err := s.genCode.GeneratePassword(length)

	if err != nil {
		return "", err
	}

	err = s.urlStorage.StoreUrl(ctx, code, url, expire)

	if err != nil {
		return "", err
	}

	return code, nil
}

func (s *shortenUrlService) GetUrlByCode(ctx context.Context, code string) (string, error) {
	return "", nil
}
