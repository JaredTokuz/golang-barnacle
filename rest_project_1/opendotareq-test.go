package main

import (
	"io/ioutil"

	"log"
	"net/http"

	"bytes"
	"encoding/json"
	// "go.mongodb.org/mongo-driver/bson/primitive"

	"time"
	"os"

	"fmt"
)

// type Match struct {
// 	Id primitive.ObjectID `bson:"_id,omitempty"`
// 	MatchId primitive.ObjectID `json:"match_id" bson:"match_id,omitempty"`
// 	TeamA []int8 `json:"teama" bson:"teama,omitempty"`
// 	TeamB []int8 `json:"teamb" bson:"teamb,omitempty"`
// 	TeamAWin bool `json:"teamawin" bson:"teamawin,omitempty"`
// 	StartTime int64 `json:"start_time" bson:"start_time,omitempty"`
// }

// type Match struct {
// 	MatchId primitive.ObjectID `bson:"match_id,omitempty"`
// 	TeamA []int8 `bson:"teama,omitempty"`
// 	TeamB []int8 `bson:"teamb,omitempty"`
// 	TeamAWin bool `bson:"teamawin,omitempty"`
// 	StartTime int64 `bson:"start_time,omitempty"`
// }

type Match struct {
	MatchId int `json:"match_id,omitempty"`
	TeamA []int8 `json:"teama,omitempty"`
	TeamB []int8 `json:"teamb,omitempty"`
	TeamAWin bool `json:"teamawin,omitempty"`
	StartTime int64 `json:"start_time,omitempty"`
}

type Matches []struct {
	MatchID int64 `json:"match_id,omitempty"`
	Teama []int `json:"teama,omitempty"`
	Teamb []int `json:"teamb,omitempty"`
	Teamawin bool `json:"teamawin,omitempty"`
	StartTime int64 `json:"start_time,omitempty"`
}

func main() {
	f, err := os.OpenFile("opendotareq-test.log", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)
	GetExample()
	GetExampleJsonDecode()
	// GetExampleCustomReq()
	// PostExample()
	// CustomRequest()
}

func GetExample() {
	var base_url string = "https://api.opendota.com/api/findMatches"

	resp, err := http.Get(base_url)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(string(body))

	var matches []Match
	json.Unmarshal(body, &matches)
	log.Println(matches)

	var data []map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(data)
	log.Println(data[0])
	log.Println(data[0]["match_id"])

	var matcherino interface{}
	json.NewDecoder(resp.Body).Decode(&matcherino)
	log.Println(matcherino)

	m := Matches{}
	json.NewDecoder(resp.Body).Decode(&m) // resp.Body is empty after the first consume with ioutil.ReadAll
	fmt.Printf("%+v\n", m)
}

func GetExampleJsonDecode() {
	var base_url string = "https://api.opendota.com/api/findMatches"

	resp, err := http.Get(base_url)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	m := Matches{}
	json.NewDecoder(resp.Body).Decode(&m)
	fmt.Printf("%+v\n", m)
}

func GetExampleCustomReq() {
	var base_url string = "https://api.opendota.com/api/findMatches"

	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	req, err := http.NewRequest("GET", base_url, nil)
	req.Header.Add("Accept", "application/json")
	if err != nil {
		log.Fatal(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(string(body))

}

func PostExample() {
	requestBody, err := json.Marshal(map[string]string{
		"name": "as982k22jk",
		"email": "masnun@gmail.com",
	})
	if err != nil {
		log.Fatalln(err)
	}

	resp, err := http.Post("https://httpbin.org/post", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		log.Fatalln(err)
	}
	
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(string(body))
}

func CustomRequest() {
	requestBody, err := json.Marshal(map[string]string{
		"name": "as982k22jk",
		"email": "masnun@gmail.com",
	})
	if err != nil {
		log.Fatalln(err)
	}

	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	request, err := http.NewRequest("POST", "https://httpbin.org/post", bytes.NewBuffer(requestBody))
	request.Header.Set("Content-type", "application/json")
	if err != nil {
		log.Fatalln(err)
	}

	resp, err := client.Do(request)
	if err != nil {
		log.Fatalln(err)
	}
	
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(string(body))
}
