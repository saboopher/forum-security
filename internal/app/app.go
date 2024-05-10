package app

import (
	"fmt"
	"forum/internal/background"
	"forum/internal/config"
	"forum/internal/database"
	handler "forum/internal/handlers"
	repository "forum/internal/repositories"
	service "forum/internal/services"
	"log"
	"net/http"
)

func Run() {
	config, err := config.NewConfig()
	if err != nil {
		log.Fatal(err) // handle errors properly
	}

	db, err := database.CreateDb(config)
	if err != nil {
		log.Fatal(err) // handle errors properly
	}
	err = database.InsertInitialData(config)
	if err != nil {
		log.Fatal(err)
	}

	err = database.RemoveSessions(config)
	if err != nil {
		log.Fatal(err)
	}
	
	repo := repository.NewRepository(db)
	service := service.NewService(repo)
	handler := handler.NewHandler(service)

	go background.WorkerScanBD(db)

	server := &http.Server{
		Addr:    config.Port,
		Handler: handler.Routes(),
	}

	fmt.Printf("Starting server on http://localhost%s", config.Port)

	if err = server.ListenAndServe(); err != nil {
		log.Fatal(err) // handle errors properly
	}
}
