package redirect

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/letabilis/desafio-url-shortener/internal/types"
)

type handler struct {
	svc types.RedirectService
}

func NewHandler(svc types.RedirectService) *handler {
	return &handler{
		svc: svc,
	}
}

func (h *handler) RegisterRoutes(r chi.Router) {
	r.Get("/{slug}", h.ResolveURL())
}

// ResolveURL godoc
// @Summary      Resolves a short URL to its corresponding long URL
// @Description  Fetches the RedisDB and redirects client to original URL
// @Tags         redirect
// @Param 	 slug path string true "Short URL Slug"
// @Success      302 "Redirect to the long URL"
// @Failure      404 "Short URL not found."
// @Router       /{slug} [get]
func (h *handler) ResolveURL() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slug := chi.URLParam(r, "slug")

		longURL, err := h.svc.GetLongURL(r.Context(), slug)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Redirect(w, r, longURL, http.StatusFound)
	}
}
