package main

import (
	"avito_go/pkg/handler"
	"avito_go/pkg/repository"
	"avito_go/pkg/service"
	"avito_go/pkg/shop"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	db, err := repository.NewPostgresDB(repository.Config{
		Host:     os.Getenv("DATABASE_HOST"),
		Port:     os.Getenv("DATABASE_PORT"),
		Username: os.Getenv("DATABASE_USER"),
		Password: os.Getenv("DATABASE_PASSWORD"),
		DBName:   os.Getenv("DATABASE_NAME"),
		SSLMode:  "disable",
	})
	if err != nil {
		log.Fatalf("faild connect to db: %s", err.Error())
	}
	defer db.Close()
	repository := repository.NewRepository(db)
	service := service.NewService(repository)
	handler := handler.NewHandler(service)
	server := new(shop.Server)
	if err := server.Run(os.Getenv("SERVER_PORT"), handler.InitRoutes()); err != nil {
		log.Fatalf("error while running server: %s", err.Error())
	}

}
