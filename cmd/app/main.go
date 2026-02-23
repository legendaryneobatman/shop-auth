package main

import (
	"context"
	"os"
	"os/signal"
	"shop-auth/internal/app"
	"shop-auth/internal/services/auth"
	"syscall"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func main() {
	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error loading env variables: %s", err.Error())
	}

	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})

	db, err := sqlx.Open("pgx", os.Getenv("DB_URL"))
	if err != nil {
		logrus.Fatalf("failed to connect to db: %s", err.Error())
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		logrus.Fatalf("failed to ping db: %s", err.Error())
	}

	authRepository := auth.NewRepository(db)
	authService := auth.NewService(authRepository)
	authHandler := auth.NewHandler(authService)

	port := os.Getenv("APP_HOST_PORT")

	application := app.NewApp(log, authHandler, port)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit
		log.Info("Shutting down...")
		cancel()
	}()

	if err := application.Run(ctx); err != nil && err != context.Canceled {
		log.Fatalf("server error: %s")
	}

	log.Info("Server stopped gracefully")
}
