package main

import (
	"fmt"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type Tag struct {
	ID int `json:"id"`
	Name string `json:"name"`
}

func main() {
	fmt.Println("Go MySQL Tutorial")

	db, err := sql.Open("mysql", "username:password@tcp(127.0.0.1:3306)/test")

	if err != nill {
		panic(err.Error())
	}

	defer db.Close()

	insert, err := db.Query("INSERT INTO test VALUES ( 2, 'TEST' )")

	if err != nil {
		panic(err.Error())
	}

	defer insert.Close()

	results, err := db.Query("SELECT id, name from tags")
	if err != nil {
		panic(err.Error())
	}

	for results.Next() {
		var tag Tag

		err = results.Scan(&tag.ID, &tag.Name)
		if err != nil {
			panic(err.Error())
		}

		log.Printf(tag.Name)
	}

	var tag Tag

	err = db.QueryRow("SELECT id, name FROM tags where id = ?", 2).Scan(&tag.ID, &tag.Name)
	if err != nil {
		panic(err.Error())
	}

	log.Println(tag.ID)
	log.Println(tag.Name)
}