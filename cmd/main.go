package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/letabilis/desafio-url-shortener/cmd/api"
	_ "github.com/letabilis/desafio-url-shortener/docs"
	"github.com/letabilis/desafio-url-shortener/internal/shorten"
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
	SERVER_ADDR, ok1 := os.LookupEnv("SERVER_ADDR")
	REDIS_ADDR, ok2 := os.LookupEnv("REDIS_ADDR")

	if !ok1 || !ok2 {
		slog.Error("missing environment variables", "SERVER_ADDR", SERVER_ADDR, "REDIS_ADDR", REDIS_ADDR)
		return
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     REDIS_ADDR,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	svc := shorten.NewService(rdb)

	api := api.NewAPI(svc)

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:7777/swagger/doc.json"),
	))

	r.Post("/shorten-url", api.ShortenURL())
	r.Get("/{slug}", api.ResolveURL())

	slog.Info("server listening on :7777")
	err := http.ListenAndServe(SERVER_ADDR, r)

	if err != nil {
		slog.Error("failed to serve", "error", err)
		return
	}

}
