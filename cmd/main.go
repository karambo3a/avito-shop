package main

import (
	"avito_go/shop"
	"avito_go/handler"
	"avito_go/repository"
	"avito_go/service"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	db, err := repository.NewPostgresDB(repository.Config{
		Host:     "localhost",
		Port:     "5432",
		Username: "postgres",
		Password: "postgres",
		DBName:   "postgres",
		SSLMode:  "disable",
	})
	if err != nil {
        log.Fatal(err)
    }
    defer db.Close()
	if err != nil {
		log.Fatal("faild connect to db: %s", err.Error());
	}
	repository := repository.NewRepository(db)
	service := service.NewService(repository)
	handler := handler.NewHandler(service)
	server := new(shop.Server)
	if error := server.Run("8080", handler.InitRoutes()); error != nil {
		log.Fatalf("error while running server: %s", error.Error())
	}

}
