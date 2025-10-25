package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"chatapp/internal/httpserver"
	"chatapp/internal/storage"
)

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func main() {
	// ENV
	redisAddr := getenv("REDIS_ADDR", "127.0.0.1:6379")
	historyLimit, _ := strconv.Atoi(getenv("HISTORY_LIMIT", "500"))
	addr := getenv("ADDR", ":8080")

	// Storage
	repo, err := storage.NewRedisRepo(redisAddr)
	if err != nil {
		log.Fatalf("redis connect: %v", err)
	}
	defer repo.Close()

	// Handlers
	h := httpserver.NewHandlers(repo, historyLimit)

	// Router
	r := httpserver.NewRouter(h)

	// HTTPserver
	srv := &http.Server{
		Addr:              addr,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("listening on %s (redis=%s, historyLimit=%d)", addr, redisAddr, historyLimit)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %v", err)
	}
}
