package main

import (
	"log/slog"

	"github.com/letabilis/desafio-url-shortener/cmd/api"
	_ "github.com/letabilis/desafio-url-shortener/docs"
	"github.com/letabilis/desafio-url-shortener/internal/redirect"
	"github.com/letabilis/desafio-url-shortener/internal/shorten"
	"github.com/letabilis/desafio-url-shortener/internal/utils"
	"github.com/redis/go-redis/v9"
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
	env, err := utils.LoadEnv("REDIS_ADDR", "REDIS_PASSWORD")
	if err != nil {
		slog.Error(err.Error())
		return
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     env["REDIS_ADDR"],
		Password: env["REDIS_PASSWORD"],
		DB:       0, // use default db
	})

	shortener := shorten.NewHandler(shorten.NewService(rdb))
	redirecter := redirect.NewHandler(redirect.NewService(rdb))

	api := api.NewAPI(
		":8080",
		shortener,
		redirecter,
	)

	slog.Info("server is now listening")
	err = api.Run()

	if err != nil {
		slog.Error("failed to serve", "error", err)
		return
	}

}
