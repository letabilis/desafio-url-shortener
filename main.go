package main

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/redis/go-redis/v9"
)

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

	r.Post("/shorten-url", api.ShortenURL())
	r.Get("/{slug}", api.ResolveURL())

	slog.Info("server listening on :7777")
	err := http.ListenAndServe(":7777", r)

	if err != nil {
		slog.Error("failed to serve", "error", err)
		return
	}

}
