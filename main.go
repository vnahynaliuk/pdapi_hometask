// Package main
//
// @title Pipedrive Deals API
// @version 1.0
// @description A simple proxy to the Pipedrive API for managing deals.
// @host localhost:8080
// @BasePath /
// @schemes http
//
// @securityDefinitions.apikey ApiKeyAuth
// @in query
// @name api_token
package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	_ "github.com/vnahynaliuk/pdapi_hometask/docs" // Swagger docs
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/vnahynaliuk/pdapi_hometask/handlers"
	"github.com/vnahynaliuk/pdapi_hometask/middleware"
)

func main() {
	r := mux.NewRouter()

	// Apply middleware
	r.Use(middleware.LoggingMiddleware)
	r.Use(middleware.MetricsMiddleware)

	// Register endpoints
	r.HandleFunc("/deals", handlers.GetDeals).Methods("GET")
	r.HandleFunc("/deals", handlers.AddDeal).Methods("POST")
	r.HandleFunc("/deals/{id}", handlers.UpdateDeal).Methods("PUT")

	// Prometheus metrics endpoint
	r.Handle("/metrics", promhttp.Handler())

	// Swagger UI endpoint
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	srv := &http.Server{
		Handler:      r,
		Addr:         ":8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Println("Server started on port 8080")
	log.Fatal(srv.ListenAndServe())
}
