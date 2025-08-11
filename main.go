package main

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/letabilis/desafio-url-shortener/docs"
	"github.com/redis/go-redis/v9"
	"github.com/swaggo/http-swagger"
)

// @title URL Shortener API
// @version 1.0
// @description This is my solution to the backend-br url-shortener challenge.
// @termsOfService https://github.com/backend-br/desafios/blob/master/url-shortener/PROBLEM.md

// @contact.name Sahil V. Dowlani
// @contact.url http://github.com/letabilis/desafio-url-shortener
// @contact.email sahilvdowlani@outlook.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host letabilis.github.io
// @BasePath /url-shortener
func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	svc := NewService(rdb)

	api := NewAPI(svc)

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:7777/swagger/doc.json"), //The url pointing to API definition
	))
	r.Post("/shorten-url", api.ShortenURL())
	r.Get("/{slug}", api.ResolveURL())

	slog.Info("server listening on :7777")
	err := http.ListenAndServe(":7777", r)

	if err != nil {
		slog.Error("failed to serve", "error", err)
		return
	}

}
