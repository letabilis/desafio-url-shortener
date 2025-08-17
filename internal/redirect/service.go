package redirect

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type Service struct {
	rdb *redis.Client
}

func NewService(rdb *redis.Client) *Service {
	return &Service{rdb: rdb}
}

func (s *Service) GetLongURL(ctx context.Context, slug string) (string, error) {
	longURL, err := s.rdb.Get(ctx, slug).Result()

	if err == redis.Nil {
		return "", fmt.Errorf("short URL not found")
	} else if err != nil {
		return "", fmt.Errorf("failed to retrieve from Redis: %v", err)
	}

	return longURL, nil

}
