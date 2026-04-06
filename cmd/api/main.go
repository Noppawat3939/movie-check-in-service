package main

import (
	httpDelivery "check-in/internal/delivery/http"
	"check-in/internal/infra/postgresl"
	"check-in/internal/infra/redis"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	// load env
	loadEnv()

	// initialize database
	db, err := postgresl.NewDB()
	if err != nil {
		log.Print("failed initialize database: %w", err)
	}

	defer postgresl.Close(db)

	client, err := redis.NewClient()
	if err != nil {
		log.Print("failed initialize redis: %w", err)
	}

	defer client.Close()

	port := os.Getenv("APP_INTERNAL_PORT")

	router := httpDelivery.NewRouter(db)

	// start server
	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	go func() {
		fmt.Printf("start server port %s", port)

		if err := srv.ListenAndServe(); err != nil {
			log.Fatalf("server failed to start: %v", err)
		}
	}()

	waitForShutdown(srv)
}

func waitForShutdown(server *http.Server) {
	quit := make(chan os.Signal, 1)

	signal.Notify(
		quit,
		syscall.SIGINT,
		syscall.SIGTERM)

	<-quit
	log.Println("shutting down server...")

	if err := server.Close(); err != nil {
		log.Printf("server shutdown error %v:", err)
	}

	log.Println("server stopped")
}

func loadEnv() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("failed load env %w:", err)
	}
}
