package router

import (
	"github.com/go-chi/chi"
	"github.com/josimar16/goexpert/desafio-ratelimiter/handler"
)

func InitializeRoutes(router *chi.Mux) {
	router.Get("/api/v1/healthz", handler.HealthzHandler)
	router.Get("/api/v1/zipcode/{zipcode}", handler.ZipCodeHandler)
}
