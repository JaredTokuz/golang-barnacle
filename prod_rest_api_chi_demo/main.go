package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	//"github.com/tonyalaribe/todoapi/basestructure/features/todo"	// store groups of routes in separate packages
)

func Routes() *chi.Mux {
	router := chi.NewRouter()
	router.Use(
		render.SetContentType(render.ContentTypeJSON),	// Set content-Type headers as application/json
		middleware.Logger,								// Log API request calls
		middleware.Compress(5),						// Compress results, mostly gzipping assets and json
		middleware.RedirectSlashes,						// Redirect slashes to no slash URL versions
		middleware.Recoverer,							// Recover from panics without crashing server
	)

	router.Route("/v1", func(r chi.Router) {
		r.Mount("/api/todo", TodoRoutes())				// should be stored in a package
	})

	return router
}

func main() {
	router := Routes()

	walkFunc := func(method string, route string, handler http.Handler, middlwares ... func(http.Handler) http.Handler) error {
		log.Printf("%s %s\n", method, route)	// Walk and print out all routes
		return nil
	}
	if err := chi.Walk(router, walkFunc); err != nil {
		log.Panicf("Logging err: %s\n", err.Error())	// panic if there is an error
	}

	log.Fatal(http.ListenAndServe(":8080", router))		// Note, the port is usually from env
}

type Todo struct {
	Slug string `json:"slug"`
	Title string `json:"title"`
	Body string `json:"body"`
}

func TodoRoutes() *chi.Mux {
	router := chi.NewRouter()
	router.Get("/{todoID}", GetATodo)
	router.Delete("/{todoID}", DeleteTodo)
	router.Post("/", CreateTodo)
	router.Get("/", GetAllTodos)
	return router
}
func GetATodo(w http.ResponseWriter, r *http.Request) {
	todoID := chi.URLParam(r, "todoID")
	todos := Todo{
		Slug:  todoID,
		Title: "Hello world",
		Body:  "Heloo world from planet earth",
	}
	render.JSON(w, r, todos) // A chi router helper for serializing and returning json
}

func DeleteTodo(w http.ResponseWriter, r *http.Request) {
	response := make(map[string]string)
	response["message"] = "Deleted TODO successfully"
	render.JSON(w, r, response) // Return some demo response
}

func CreateTodo(w http.ResponseWriter, r *http.Request) {
	response := make(map[string]string)
	response["message"] = "Created TODO successfully"
	render.JSON(w, r, response) // Return some demo response
}

func GetAllTodos(w http.ResponseWriter, r *http.Request) {
	todos := []Todo{
		{
			Slug:  "slug",
			Title: "Hello world",
			Body:  "Heloo world from planet earth",
		},
	}
	render.JSON(w, r, todos) // A chi router helper for serializing and returning json
}