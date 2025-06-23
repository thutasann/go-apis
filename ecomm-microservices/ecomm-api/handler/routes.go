package handler

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
)

var r *chi.Mux

func RegisterRoutes(handler *handler) *chi.Mux {
	r = chi.NewRouter()

	r.Route("/products", func(r chi.Router) {
		r.Post("/", handler.createProduct)
	})

	fmt.Println(":::: successfully register routes ::::")

	return r
}

func Start(addr string) error {
	return http.ListenAndServe(addr, r)
}
