package main

import (
	"CurrencyClient/iternal/httpserver/handlers/allcurrency"
	"CurrencyClient/iternal/lib/logger/sl"
	"CurrencyClient/storage/postgres"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/joho/godotenv"
	"log/slog"
	"net/http"
	"os"
	"time"
)

func main() {
	log := slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
	)

	log.Info("Starting getCurrencu ", slog.String("env", slog.LevelInfo.String()))
	log.Debug("Debug message are enabled")

	// TODO: get env
	err := godotenv.Load(".env")
	if err != nil {
		log.Error("can't get environment", sl.Err(err))
		os.Exit(1)
	}

	// TODO: init storage
	storagePath := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		os.Getenv("host"), os.Getenv("port"), os.Getenv("user"),
		os.Getenv("password"), os.Getenv("dbname"))

	storage, err := postgres.New(storagePath)

	if err != nil {
		log.Error("failed to init db", sl.Err(err))

		os.Exit(1)
	}

	// TODO: init router
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer) // востановление работы приложение при панике в любом handler
	router.Use(middleware.URLFormat) // обработка URL для вида ссылки: r.GET("/article/{id},...)}

	// load dates from api
	go func() {
		for {

			var request []allcurrency.Request

			res, err := http.Get("https://api.coingecko.com/api/v3/coins/markets?vs_currency=usd&order=market_cap_desc&per_page=250&page=1")

			if err != nil {
				log.Error("failed to get currency")

				os.Exit(1)
			}
			defer res.Body.Close()

			err = render.DecodeJSON(res.Body, &request)

			if err != nil {
				log.Error("Failed to decode request", sl.Err(err))

				os.Exit(1)
			}

			err = storage.Save(request)
			if err != nil {
				log.Error("failed to save dates in DB")

				os.Exit(1)
			}
			time.Sleep(10 * time.Minute)
		}
	}()

	router.Get("/Get/", allcurrency.New(log, storage))

	// TODO: Start server

	log.Info("Starting server",
		slog.String("addres", "localhost"),
	)

	server := http.Server{
		Addr:         "localhost:8080",
		Handler:      router,
		ReadTimeout:  time.Duration(4),
		WriteTimeout: time.Duration(4),
		IdleTimeout:  time.Duration(60),
	}

	if err := server.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}

	log.Error("server stop")

}
