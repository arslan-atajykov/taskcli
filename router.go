package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func NewRouter() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to the task manager api"))
	})

	r.Route("/tasks", func(r chi.Router) {
		r.Get("/", GetAllTasks)
		r.Post("/", CreateTask)
		r.Get("/{id}", GetTaskByID)
		r.Delete("/{id}", DeleteTask)
		r.Put("/{id}", UpdateTask)
		r.Get("/", GetAllFilter)
	})
	return r

}
