package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"service/internal/cache"
	"service/internal/config"
	"service/internal/http-server/handlers"
	"service/internal/storage/postgres"
	"service/natsStreaming"
)

func main() {
	// init config: cleanenv
	cfg := InitConfig()

	// init sync
	var wg sync.WaitGroup
	defer wg.Wait()

	// init storage: postgres
	storage := InitStorage(cfg)

	// init router: chi, chi render
	router := chi.NewRouter()

	// init middleware: chi Mux, middleware
	InitMiddleware(router)

	// init cache go-cache
	cache := InitCache(storage, &wg)

	// init handlers
	InitHandlers(router, storage, cache)

	// init nats-streaming
	InitNatsStreaming(&wg, storage, cache)

	// run server
	RunServer(cfg, router)
}

func InitConfig() *config.Config {
	configPath := flag.String("CONFIG_PATH", "", "path to config")
	flag.Parse()
	return config.MustLoad(*configPath)
}

func InitStorage(cfg *config.Config) *postgres.Storage {
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

func InitMiddleware(router chi.Router) {
	router.Use(middleware.RequestID) //
	router.Use(middleware.Logger)    //
	router.Use(middleware.Recoverer) // обработка panic в handler
	router.Use(middleware.URLFormat) // обработка url
}

func InitCache(storage *postgres.Storage, wg *sync.WaitGroup) *cache.Cache {
	cache, err := cache.New(storage, wg)
	if err != nil {
		log.Println("Can not init cache")
	}
	return cache
}

func InitHandlers(router *chi.Mux, storage *postgres.Storage, cache *cache.Cache) {
	// HOME
	router.Get("/", handlers.HomePage)

	router.Get("/add", handlers.AddOrderPage(""))
	router.Post("/add", handlers.AddOrder(storage, cache))

	// FIND
	router.Get("/find", handlers.FindOrderByIDPage(""))
	router.Post("/find", handlers.FindOrderByID(storage, cache))

	// ORDER DETAILS
	router.Get("/order", handlers.OrderDetailsPage(nil))
}

func InitNatsStreaming(wg *sync.WaitGroup, storage *postgres.Storage, cache *cache.Cache) {
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer fmt.Println("Shutting down...")
		defer wg.Done()

		if err := natsStreaming.RunNatsStreaming(storage, cache); err != nil {
			log.Println(err)
		}
	}(wg)
}

func RunServer(cfg *config.Config, router *chi.Mux) {
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
