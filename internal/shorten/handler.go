package shorten

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"
	"github.com/letabilis/desafio-url-shortener/internal/types"
	"github.com/letabilis/desafio-url-shortener/internal/utils"
)

type handler struct {
	svc types.ShortenService
}

func NewHandler(svc types.ShortenService) *handler {
	return &handler{
		svc: svc,
	}
}

func (h *handler) RegisterRoutes(r chi.Router) {
	r.Post("/shorten-url", h.ShortenURL())
	r.Get("/{slug}", h.ResolveURL())
}

// ShortenURL godoc
// @Summary      Shorten a Long URL
// @Description  Responds with a Slug (shortCode) and an Expiry Date (1 day by default)
// @Tags         url
// @Param request body types.AnyRequest true "The URL to shorten"
// @Accept       json
// @Produce      json
// @Success      200 {object} types.ShortenResponse
// @Failure      400 {string} string "Bad Request"
// @Failure      500 {string} string "Internal Server Error"
// @Router       /shorten-url [post]
func (h *handler) ShortenURL() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var payload types.AnyRequest
		err := json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		url, err := url.ParseRequestURI(payload.URL)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		shortenResponse, err := h.svc.GetSlug(r.Context(), url.String())

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		utils.WriteJSON(w, http.StatusOK, shortenResponse)
	}
}

// ResolveURL godoc
// @Summary      Resolves a short URL to its corresponding long URL
// @Description  Fetches the RedisDB and redirects client to original URL
// @Tags         url
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
