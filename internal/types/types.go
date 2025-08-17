package types

import (
	"context"
	"time"

	"github.com/go-chi/chi/v5"
)

// The AnyRequest struct maps incoming Shortening or Redirection JSON requests.
type AnyRequest struct {
	URL string `json:"url" example:"https://example.com"`
}

// The ShortenResponse struct consists of the resulting output of a ShortenRequest.
type ShortenResponse struct {
	Slug   string    `json:"slug" example:"P3Iww4CcYhA"`
	Expiry time.Time `json:"expiry" example:"2025-08-12T12:19:22.040906963-04:00"`
}

type Handler interface {
	RegisterRoutes(r chi.Router)
}

type ShortenService interface {
	GetSlug(ctx context.Context, longURL string) (*ShortenResponse, error)
}

type RedirectService interface {
	GetLongURL(ctx context.Context, slug string) (string, error)
}
