package main

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"labkoding.my.id/kasir-api/handler"
)

func handleError(w http.ResponseWriter, err error) {
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	r.Get("/categories", func(w http.ResponseWriter, r *http.Request) {
		err := handler.GetAll(w, r)
		handleError(w, err)
	})

	r.Post("/categories", func(w http.ResponseWriter, r *http.Request) {
		err := handler.CreateCategory(w, r)
		handleError(w, err)
	})

	r.Get("/categories/{id}", func(w http.ResponseWriter, r *http.Request) {
		err := handler.FindCategoryById(w, r)
		handleError(w, err)
	})

	r.Put("/categories/{id}", func(w http.ResponseWriter, r *http.Request) {
		err := handler.UpdateCategoryById(w, r)
		handleError(w, err)
	})

	r.Delete("/categories/{id}", func(w http.ResponseWriter, r *http.Request) {
		err := handler.DeleteCategoryById(w, r)
		handleError(w, err)
	})

	http.ListenAndServe(":3000", r)
}
