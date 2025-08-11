package main

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type API struct {
	svc *Service
}

func NewAPI(svc *Service) *API {
	return &API{
		svc: svc,
	}
}

// ShortenURL godoc
// @Summary      Shorten a Long URL
// @Description  Responds with a Slug (shortCode) and an Expiry Date (1 day by default)
// @Tags         url
// @Param request body AnyRequest true "The URL to shorten"
// @Accept       json
// @Produce      json
// @Success      200 {object} ShortenResponse
// @Failure      400 {string} string "Bad Request"
// @Failure      500 {string} string "Internal Server Error"
// @Router       /shorten [post]
func (api *API) ShortenURL() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var long AnyRequest
		err := json.NewDecoder(r.Body).Decode(&long)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		shortenResponse, err := api.svc.GetSlug(r.Context(), long.URL)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		WriteJSON(w, http.StatusOK, shortenResponse)
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
func (api *API) ResolveURL() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slug := chi.URLParam(r, "slug")
		longURL, err := api.svc.GetLongURL(r.Context(), slug)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Redirect(w, r, longURL, http.StatusFound)
	}
}
