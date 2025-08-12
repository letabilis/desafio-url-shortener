package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/letabilis/desafio-url-shortener/internal/types"
	"github.com/letabilis/desafio-url-shortener/internal/utils"

	"github.com/swaggo/http-swagger"
)

type API struct {
	addr      string
	shortener types.ShortenService
}

func NewAPI(addr string, shortener types.ShortenService) *API {
	return &API{
		addr:      addr,
		shortener: shortener,
	}
}

func (api *API) mount() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:7777/swagger/doc.json"),
	))

	r.Post("/shorten-url", api.ShortenURL())
	r.Get("/{slug}", api.ResolveURL())

	return r
}

func (api *API) Run() error {
	return http.ListenAndServe(api.addr, api.mount())
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
func (api *API) ShortenURL() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var long types.AnyRequest
		err := json.NewDecoder(r.Body).Decode(&long)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		shortenResponse, err := api.shortener.GetSlug(r.Context(), long.URL.String())

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
func (api *API) ResolveURL() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slug := chi.URLParam(r, "slug")
		longURL, err := api.shortener.GetLongURL(r.Context(), slug)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Redirect(w, r, longURL, http.StatusFound)
	}
}
