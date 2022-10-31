package main

import (
	// "os"
	"fmt"
	"net/http"
	"encoding/json"
	"time"
)

type Matches []struct {
	MatchID int64 `json:"match_id"`
	Teama []int `json:"teama"`
	Teamb []int `json:"teamb"`
	Teamawin bool `json:"teamawin"`
	StartTime int64 `json:"start_time"`
}

var myClient = &http.Client{Timeout: 10 * time.Second}

func getJson(url string, target interface{}) error {
	r, err := myClient.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
	
}

var api_url = "https://api.opendota.com/api/findMatches"

func main() {
	foo1 := new(Matches)
	getJson(api_url, foo1)
	println(foo1)

	foo2 := Matches{}
	getJson(api_url, &foo2)
	println(foo2[0].MatchID)

	fmt.Printf("%+v\n", foo2)

	var foo3 []interface{}
	getJson(api_url, &foo3)
	fmt.Printf("%+v\n", foo3)
	fmt.Printf("%+v\n", foo3[0])
	fmt.Printf("%+v\n", foo3[0])

	foo4 := Matches{}
	getJson(api_url, &foo4)
	m := []interface{}{foo4}
	fmt.Printf("\n\n")
	fmt.Printf("%+v\n", m)
}