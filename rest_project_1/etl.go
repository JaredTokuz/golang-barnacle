package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)
// break up the code by splitting route groups in their own files.
type dotaETLResource struct{}

// Full URL example: "https://api.opendota.com/api/matches/18029376?api_key=YOUR_API_KEY"
var base_url string = "https://api.opendota.com/api/"

// Routes creates a REST router for the todos resource
func (rs dotaETLResource) Routes(api_key string) chi.Router {
	r := chi.NewRouter()
	// r.Use() // some middleware..

	r.Get("/", rs.GetRandomSample)

	return r
}

func (rs dotaETLResource) GetRandomSample(w http.ResponseWriter, r *http.Request) {

}