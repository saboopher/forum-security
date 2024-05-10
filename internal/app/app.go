package app

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"

	"forum/internal/background"
	"forum/internal/config"
	"forum/internal/database"
	handler "forum/internal/handlers"
	repository "forum/internal/repositories"
	service "forum/internal/services"
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

	tlsConf := tls.Config{
		PreferServerCipherSuites: true,
		CurvePreferences:         []tls.CurveID{tls.CurveP384, tls.CurveP256},
	}

	server := &http.Server{
		Addr:      config.Port,
		Handler:   handler.Routes(),
		TLSConfig: &tlsConf,
	}

	fmt.Printf("Starting server on http://localhost%s", config.Port)

	// if err = server.ListenAndServe(); err != nil {
	// 	log.Fatal(err) // handle errors properly
	// }
	fmt.Println(config.CertTLS, config.KeyTLS)

	if err = server.ListenAndServeTLS(config.CertTLS, config.KeyTLS); err != nil {
		log.Fatal("ListenAndServeTLS: ", err)
	}
}
