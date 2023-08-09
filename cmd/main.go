package main

import (
	"flag"
	"fmt"
	"github.com/go-chi/chi/middleware"
	"log"
	"net/http"

	"L0/internal/config"
	"L0/internal/storage/postgres"
	"github.com/go-chi/chi"
)

func main() {
	// init config: cleanenv
	cfg := initConfig()

	// init storage: postgres
	storage := initStorage(cfg)

	// init router: chi, chi render
	router := chi.NewRouter()

	// init middleware: chi Mux, middleware
	initMiddleware(router)

	// init handlers
	initHandlers(cfg, router, storage)

	// run server
	runServer(cfg, router)
}

func initConfig() *config.Config {
	configPath := flag.String("CONFIG_PATH", "", "path to config")
	flag.Parse()
	return config.MustLoad(*configPath)
}

func initStorage(cfg *config.Config) *postgres.Storage {
	connectionString := fmt.Sprintf(
		"host=%s port=%d user=%s dbname=%s sslmode=%s",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.User,
		cfg.Database.DBName, cfg.Database.SSLMode,
	)

	storage, err := postgres.New(connectionString)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to init storage: %s", err))
	}

	return storage
}

func initMiddleware(router chi.Router) {
	router.Use(middleware.RequestID) //
	router.Use(middleware.Logger)    //
	router.Use(middleware.Recoverer) // обработка panic в handler
	router.Use(middleware.URLFormat) // обработка url
}

func initHandlers(cfg *config.Config, router *chi.Mux, storage *postgres.Storage) {
	// TODO: handler
	// router.Get("/{id}")
}

func runServer(cfg *config.Config, router *chi.Mux) {
	log.Printf("starting server\naddress: %s", cfg.HttpServer.Address)

	srv := &http.Server{
		Addr:         cfg.HttpServer.Address,
		Handler:      router,
		ReadTimeout:  cfg.HttpServer.Timeout,
		WriteTimeout: cfg.HttpServer.Timeout,
		IdleTimeout:  cfg.HttpServer.IdleTimeout,
	}
	if err := srv.ListenAndServe(); err != nil {
		log.Println("failed to start server")
	}

	log.Println("server stopped")
}