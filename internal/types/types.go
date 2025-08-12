package types

import (
	"context"
	"net/url"
	"time"
)

// The AnyRequest struct maps incoming Shortening or Redirection JSON requests.
type AnyRequest struct {
	URL url.URL `json:"url" example:"https://example.com"`
}

// The ShortenResponse struct consists of the resulting output of a ShortenRequest.
type ShortenResponse struct {
	Slug   string    `json:"slug" example:"P3Iww4CcYhA"`
	Expiry time.Time `json:"expiry" example:"2025-08-12T12:19:22.040906963-04:00"`
}

type ShortenService interface {
	GetSlug(ctx context.Context, longURL string) (ShortenResponse, error)
	GetLongURL(ctx context.Context, slug string) (string, error)
}
