package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	apirouter "github.com/muosilva/url-shortener/cmd/router"
	"github.com/muosilva/url-shortener/internal/handler"
	"github.com/muosilva/url-shortener/internal/service"
	"github.com/muosilva/url-shortener/storage/postgres"
	rediscache "github.com/muosilva/url-shortener/storage/redis"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/redis/go-redis/v9"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db, err := openPostgres(ctx)
	if err != nil {
		return err
	}
	defer db.Close()

	redisClient := openRedis()
	defer redisClient.Close()

	if err := redisClient.Ping(ctx).Err(); err != nil {
		return err
	}

	urlRepository := postgres.NewPostgresURLRepository(db)
	urlCache := rediscache.NewRedisURLCache(redisClient)
	urlService := service.NewService(urlRepository, urlCache)
	urlHandler := handler.NewHandler(urlService)

	server := &http.Server{
		Addr:         getEnv("HTTP_ADDR", ":8080"),
		Handler:      apirouter.New(urlHandler),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("listening on %s", server.Addr)
	return server.ListenAndServe()
}

func openPostgres(ctx context.Context) (*sql.DB, error) {
	db, err := sql.Open("pgx", getEnv("DATABASE_URL", "postgres://url_shortener:url_shortener@localhost:5432/url_shortener?sslmode=disable"))
	if err != nil {
		return nil, err
	}

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

func openRedis() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: getEnv("REDIS_ADDR", "localhost:6379"),
	})
}

func getEnv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}
