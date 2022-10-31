package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type usersResource struct{}

func (rs usersResource) Routes() chi.Router {
	r := chi.NewRouter()
	// r.Use()	// some middleware..

	r.Get("/", rs.List)	// GET /users -read a list of users
	r.Post("/", rs.Create)	// POST /users -create a new user and persist it
	r.Delete("/", rs.Delete)	// DELETE /users -read a list of users

	r.Route("/{id}", func(r chi.Router) {
		// r.Use(rs.TodoCtx)	// lets have a todos map, and lets actually load/manipulate
		r.Get("/", rs.Get)		// GET /todos/{id} - read a single user by :id
		r.Put("/", rs.Update)	// PUT /todos/{id} - update a single user by :id
		r.Delete("/", rs.Delete)	// DELETE /todos/{id} - delete a single user by :id
	})

	return r
}

func (rs usersResource) List(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("todos list of stuff.."))
}

func (rs usersResource) Create(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("todos create"))
}

func (rs usersResource) Get(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("todo get"))
}

func (rs usersResource) Update(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("todo update"))
}

func (rs usersResource) Delete(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("todo delete"))
}