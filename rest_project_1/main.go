package main

import (
	"context"
	"fmt"
	//"os"
	"time"
	"log"
	"crypto/tls"

	// "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	// "go.mongodb.org/mongo-driver/mongo/readpref"
	
	"github.com/joho/godotenv"

	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"encoding/json"
	"io/ioutil"
	//"bytes"
)

var dbclient *mongo.Client
var myEnv map[string]string

func main() {
	// Load in env file
	myEnv, err := godotenv.Read()
	if err != nil {
		log.Fatal("Error load .env file")
	}

	// Set up the database client connection
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	dbclient, err := mongo.Connect(ctx, options.Client().ApplyURI(myEnv["ATLAS_URI"]))
	if err != nil {
		panic(err)
	}
	defer dbclient.Disconnect(ctx)
	// Test the connection
	// err = client.Ping(ctx, readpref.Primary())
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// Database and collections
	// database := client.Database("quickstart")
	// matchesCollection := database.Collection("matches")

	// Set up chi router
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("The proof of concept REST api!"))
	})

	r.Mount("/dota-etl", dotaETLResource{}.Routes())
	//r.Mount("/dota-games", dotaResources{}.Routes())

	// generate a `Certificate` struct
	cert, _ := tls.LoadX509KeyPair("localhost.crt", "localhost.key")
	// create custom server
	s := &http.Server{
		Addr: ":8888",
		Handler: r,
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{ cert },
		},
	}

	s.ListenAndServeTLS("", "")
}

// Dota ETL Routes .. put this in its own file for modulated code.
type dotaETLResource struct{
	DB_URI string
	API_KEY string
}

// Full URL example: "https://api.opendota.com/api/matches/18029376?api_key=YOUR_API_KEY"
var base_url string = "https://api.opendota.com/api"

// Routes creates a REST router for the todos resource
func (rs dotaETLResource) Routes() chi.Router {
	r := chi.NewRouter()
	// r.Use() // some middleware..
	//rs.DB_URI = env["ATLAS_URI"]

	r.Get("/", rs.GetRandomSample)

	return r
}

type Match struct {
	Id primitive.ObjectID `bson:"_id,omitempty"`
	MatchId primitive.ObjectID `json:"match_id"` `bson:"match_id,omitempty"`
	TeamA []int8 `json:"teama"` `bson:"teama,omitempty"`
	TeamB []int8 `json:"teamb"` `bson:"teamb,omitempty"`
	TeamAWin bool `json:"teamawin"` `bson:"teamawin,omitempty"`
	StartTime int64 `json:"start_time"` `bson:"start_time,omitempty"`
}

type Podcast struct {
	ID     primitive.ObjectID `bson:"_id,omitempty"`
	Title  string             `bson:"title,omitempty"`
	Author string             `bson:"author,omitempty"`
	Tags   []string           `bson:"tags,omitempty"`
}

func (rs dotaETLResource) GetRandomSample(w http.ResponseWriter, r *http.Request) {
	req_url := fmt.Sprintf("%s/findMatches", base_url)

	// var req_url string = "https://api.opendota.com/api/findMatches"

	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	req, err := http.NewRequest("GET", req_url, nil)
	req.Header.Add("Accept", "application/json")
	if err != nil {
		log.Fatal(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	// var result map[string]interface{}
	// json.NewDecoder(resp.Body).Decode(&result)
	// log.Println(result)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	var matches []Match
	json.Unmarshal(body, &matches)
	fmt.Println(matches)
	
	// ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	// // db_uri := rs.DB_URI
	// dbclient, err := mongo.Connect(ctx, options.Client().ApplyURI(rs.DB_URI))
	// if err != nil {
	// 	panic(err)
	// }
	// defer dbclient.Disconnect(ctx)

	// Database and collections
	database := dbclient.Database("dota")
	matchesCollection := database.Collection("matches")
	
	podcast := []interface{}{
		Podcast{
			Title:  "The Polyglot Developer",
			Author: "Nic Raboy",
			Tags:   []string{"development", "programming", "coding"},
		},
		Podcast{
			Title:  "The Polyglot Magi",
			Author: "Go Babox",
			Tags:   []string{"development", "programming", "coding"},
		},
	}

	// Insert the requested records
	insertResult, err := matchesCollection.InsertMany(ctx, podcast)
	if err != nil {
		panic(err)
	}
	fmt.Println(insertResult)

	w.Write([]byte(body))
}
