package shorten

import (
	"context"
	"log/slog"
	"time"

	"github.com/letabilis/desafio-url-shortener/internal/types"
	"github.com/redis/go-redis/v9"
)

type Service struct {
	rdb *redis.Client
}

func NewService(rdb *redis.Client) *Service {
	return &Service{rdb: rdb}
}

func (s *Service) GetSlug(ctx context.Context, longURL string) (*types.ShortenResponse, error) {
	slug := GetShortCode(longURL)
	expiry := 24 * time.Hour
	err := s.rdb.Set(ctx, slug, longURL, expiry).Err()
	if err != nil {
		slog.Error("unable to set slug", "error", err)
		return nil, err
	}
	return &types.ShortenResponse{Slug: slug, Expiry: time.Now().Add(expiry)}, nil
}
