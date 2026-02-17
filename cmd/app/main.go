package main

import (
	"os"
	shop "shop-auth"
	"shop-auth/internal/bootstrap"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error loading env variables: %s", err.Error())
	}

	db, err := sqlx.Open("pgx", os.Getenv("SHARED_DB_URL"))

	if err != nil {
		logrus.Fatalf("failed to connect to db %s", err.Error())
	}

	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:    []string{"Origin", "Content-Type", "Accept", "Authorization"},
	}))
	cors.Default()
	routesHandler := bootstrap.NewHandler(db)
	routesHandler.Init(router)

	server := new(shop.Server)
	if err := server.Run(os.Getenv("APP_HOST_PORT"), router); err != nil {
		logrus.Fatalf("error occurred while running http server: %s", err.Error())
	} else {
		logrus.Infoln("Server started on port: ", os.Getenv("APP_HOST_PORT"))
	}
}
