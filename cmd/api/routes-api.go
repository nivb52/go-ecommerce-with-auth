package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"}, //FOR DEV ENV
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Contnet-Type", "X-CSRF-Token"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	mux.Get("/api/liveness", app.Liveness)
	mux.Get("/api/liveness/", app.Liveness)
	mux.Post("/api/payment-intent", app.GetPaymentIntent)
	mux.Post("/api/payment-intent/", app.GetPaymentIntent)

	return mux
}
