package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	"github.com/gorilla/mux"
	// "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"github.com/joho/godotenv"
	"io/ioutil"
	"os"
)

var dbclient *mongo.Client
var base_url string = "https://api.opendota.com/api"
var myEnv map[string]string
type Person struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Firstname string             `json:"firstname,omitempty" bson:"firstname,omitempty"`
	Lastname  string             `json:"lastname,omitempty" bson:"lastname,omitempty"`
}
type Matches []struct {
	MatchID int64 `json:"match_id,omitempty"`
	Teama []int `json:"teama,omitempty"`
	Teamb []int `json:"teamb,omitempty"`
	Teamawin bool `json:"teamawin,omitempty"`
	StartTime int64 `json:"start_time,omitempty"`
}
func CreatePersonEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var person Person
	_ = json.NewDecoder(request.Body).Decode(&person)
	collection := dbclient.Database("dota").Collection("people")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, _ := collection.InsertOne(ctx, person)
	json.NewEncoder(response).Encode(result)
}
func GetPeopleEndpoint(response http.ResponseWriter, request *http.Request) {
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
	// body, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	// log.Println(string(body))
	// var matches []Match
	// json.Unmarshal(body, &matches)
	// log.Println(matches)

	var matches Matches{}
	json.NewDecoder(resp.Body).Decode(&matches)
	fmt.Printf("%+v\n", matches)

	response.Header().Set("content-type", "application/json")

	collection := dbclient.Database("dota").Collection("matches")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)

	// var insertion []interface{}
	// for _, t := range matches{
	// 	insertion = append(insertion, t)
	// }

	insertResult, err := collection.InsertMany(ctx, matches)
	if err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(response).Encode(insertResult)

	// if err != nil {
	// 	response.WriteHeader(http.StatusInternalServerError)
	// 	response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
	// 	return
	// }
	// defer cursor.Close(ctx)
	// for cursor.Next(ctx) {
	// 	var person Person
	// 	cursor.Decode(&person)
	// 	people = append(people, person)
	// }
	// if err := cursor.Err(); err != nil {
	// 	response.WriteHeader(http.StatusInternalServerError)
	// 	response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
	// 	return
	// }
	// json.NewEncoder(response).Encode(people)
}
func GetPersonEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	var person Person
	collection := dbclient.Database("dota").Collection("people")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	err := collection.FindOne(ctx, Person{ID: id}).Decode(&person)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(person)
}

func main() {
	// Load in env file
	myEnv, err := godotenv.Read()
	if err != nil {
		log.Fatal("Error load .env file")
	}
	f, err := os.OpenFile("gorilla-mux-man.log", os.O_RDWR, | os.O_CREATE | os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)
	fmt.Println("Starting the application...")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	clientOptions := options.Client().ApplyURI(myEnv["ATLAS_URI"])
	dbclient, _ = mongo.Connect(ctx, clientOptions)
	defer dbclient.Disconnect(ctx)
	router := mux.NewRouter()
	router.HandleFunc("/person", CreatePersonEndpoint).Methods("POST")
	router.HandleFunc("/matches", GetPeopleEndpoint).Methods("GET")
	router.HandleFunc("/person/{id}", GetPersonEndpoint).Methods("GET")
	http.ListenAndServe(":12345", router)
}