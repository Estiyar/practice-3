package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"golang/internal/handler"
	"golang/internal/middleware"
	"golang/internal/repository"
	"golang/internal/repository/_postgres"
	"golang/internal/usecase"
	"golang/pkg/modules"
)

func main() {
	cfg := &modules.PostgreConfig{
		Host:        "localhost",
		Port:        "5432",
		Username:    "postgres",
		Password:    "Esti2005",
		DBName:      "mydb",
		SSLMode:     "disable",
		ExecTimeout: 5 * time.Second,
	}

	pg := _postgres.NewPGXDialect(context.Background(), cfg)
	repos := repository.NewRepositories(pg)
	uc := usecase.NewUserUsecase(repos.UserRepository)
	h := handler.NewUserHandler(uc)

	mux := http.NewServeMux()
	mux.HandleFunc("/health", handler.Health)
	mux.HandleFunc("/users", h.Users)
	mux.HandleFunc("/users/", h.UserByID)

	apiKey := "Estiyar"

	finalHandler := middleware.Logging(middleware.APIKey(apiKey)(mux))

	server := &http.Server{
		Addr:    ":8080",
		Handler: finalHandler,
	}

	log.Println("listening on :8080")
	log.Fatal(server.ListenAndServe())
}