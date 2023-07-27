package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*", "http://localhost:4000", "localhost:4001"}, //FOR DEV ENV
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Contnet-Type", "X-CSRF-Token", "Origin", "Content-Length"},
		AllowCredentials: false,
		MaxAge:           300,
	}))
	mux.Use(SessionLoad)

	mux.Get("/", app.Home)
	mux.Get("/liveness", app.Liveness)
	mux.Get("/virtual-terminal", app.VirtualTerminal)
	mux.Post("/payment-succeeded", app.PaymentSucceeded)
	mux.Post("/receipt", app.Receipt)

	mux.Get("/widgets/{id}", app.WidgetById)

	fileServer := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}
