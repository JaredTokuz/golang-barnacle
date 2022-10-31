package main

import (
	"log"
	"os"
	"fmt"
	"net/http"
	"encoding/json"
	"time"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson"
	"github.com/joho/godotenv"
)

type Matches []struct {
	MatchID int64 `json:"match_id"`
	Teama []int `json:"teama"`
	Teamb []int `json:"teamb"`
	Teamawin bool `json:"teamawin"`
	StartTime int64 `json:"start_time"`
}

var myClient = &http.Client{Timeout: 10 * time.Second}
var base_url = "https://api.opendota.com/api"
var api_url = "https://api.opendota.com/api/findMatches"
var dbclient *mongo.Client
var myEnv map[string]string
var logname = "opendota-etl.log"

func GetRequestJson(endpoint string, target interface{}) {
	req_url := fmt.Sprintf("%s/%s", base_url, endpoint)
	log.Printf("url: %+v\n", req_url)

	req, err := http.NewRequest("GET", api_url, nil)
	req.Header.Add("Accept", "application/json")
	if err != nil {
		log.Fatal(err)
	}
	getJson(req, &target)
	log.Printf("interface: %+v\n", target)
}

func getJson(request *http.Request, target interface{}) error {
	r, err := myClient.Do(request)
	if err != nil {
		log.Printf("Client Do error : %+v\n", target)
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
	
}

func ImplementLogs(name string) (*os.File, error) {
	f, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
		return nil, err
	}
	// defer f.Close() // on the outside
	f.Truncate(0)
	log.SetOutput(f)
	return f, nil
}

func FindMatches() {
	matches := Matches{}
	GetRequestJson("findMatches", &matches)

	var operations []mongo.WriteModel

	for _, m := range matches {
		operation_i := mongo.NewUpdateOneModel()
		operation_i.SetFilter(bson.M{"match_id": m.MatchID})
		operation_i.SetUpdate(bson.D{
			{
				"$set", bson.D{
					{"match_id", m.MatchID},
					{"teama", m.Teama},
					{"teamb", m.Teamb},
					{"teamawin", m.Teamawin},
					{"start_time", m.StartTime},
				}},
		})
		operation_i.SetUpsert(true)
		operations = append(operations, operation_i)
	}

	log.Printf("Operations: %+v\n", operations)

	myEnv, err := godotenv.Read()
	if err != nil {
		log.Fatal("Error load .env file")
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	clientOptions := options.Client().ApplyURI(myEnv["ATLAS_URI"])
	dbclient, _ = mongo.Connect(ctx, clientOptions)
	defer dbclient.Disconnect(ctx)
	collection := dbclient.Database("dota").Collection("matches")

	bulkOption := options.BulkWriteOptions{}
	bulkOption.SetOrdered(true)
	result, err := collection.BulkWrite(ctx, operations, &bulkOption)
	if err != nil {
		log.Fatalf("BulkWrite Failed: %s.", err)
	}
	log.Printf("Results: %s", result)

}
func main() {
	f, _ := ImplementLogs(logname)
	defer f.Close()
	log.SetOutput(f)
	// Further optimize add the env code out here and test it
	FindMatches()
}