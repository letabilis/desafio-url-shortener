package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/letabilis/desafio-url-shortener/internal/types"

	"github.com/swaggo/http-swagger"
)

type API struct {
	addr     string
	handlers []types.Handler
}

func NewAPI(addr string, handlers ...types.Handler) *API {
	return &API{
		addr:     addr,
		handlers: handlers,
	}
}

func (api *API) mount() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:7777/swagger/doc.json"),
	))

	for _, handler := range api.handlers {
		handler.RegisterRoutes(r)
	}

	return r
}

func (api *API) Run() error {
	return http.ListenAndServe(api.addr, api.mount())
}
