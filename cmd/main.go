package main

import (
	"log/slog"
	"os"

	"github.com/letabilis/desafio-url-shortener/cmd/api"
	_ "github.com/letabilis/desafio-url-shortener/docs"
	"github.com/letabilis/desafio-url-shortener/internal/shorten"
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

	api := api.NewAPI(SERVER_ADDR, svc)

	slog.Info("server listening on", "addr", SERVER_ADDR)
	err := api.Run()

	if err != nil {
		slog.Error("failed to serve", "error", err)
		return
	}

}
