package service

import (
	"context"
	"time"

	"github.com/NghiaLeopard/bookmark-management/internal/repository"
)

//go:generate mockery --name ShortenUrlService --filename shortenurl.go
type ShortenUrl interface {
	CreateShortenUrl(ctx context.Context, url string, expire time.Duration) (string, error)
	GetUrlByCode(ctx context.Context, code string) (string, error)
}

type shortenUrlService struct {
	urlStorage repository.UrlStorage
	genCode    GenPass
}

func NewShortenUrl(urlStorage repository.UrlStorage, genCode GenPass) ShortenUrl {
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

	_, err = s.urlStorage.GetUrlByCode(ctx, code)

	if err == nil {
		return s.CreateShortenUrl(ctx, url, expire)
	}

	err = s.urlStorage.StoreUrl(ctx, code, url, expire)

	if err != nil {
		return "", err
	}

	return code, nil
}

func (s *shortenUrlService) GetUrlByCode(ctx context.Context, code string) (string, error) {
	return s.urlStorage.GetUrlByCode(ctx, code)
}
