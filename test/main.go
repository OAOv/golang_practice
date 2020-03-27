package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	//_ "github.com/go-sql-driver/mysql"
)

type event struct {
	ID   string `json:"ID"`
	Name string `json:"Name"`
	Age  string `json:"Age"`
}

var allEvents []event

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi, %s", r.URL.Path[1:])
}

func createHandler(w http.ResponseWriter, r *http.Request) {
	ID, err1 := r.URL.Query()["ID"]
	name, err2 := r.URL.Query()["name"]
	age, err3 := r.URL.Query()["age"]
	if !err1 || !err2 || !err3 {
		errorInput()
		return
	}
	e := &event{ID[0], name[0], age[0]}
	allEvents = append(allEvents, *e)
	e_json, _ := json.Marshal(e)
	fmt.Fprintf(w, string(e_json))
}

func readAllHandler(w http.ResponseWriter, r *http.Request) {
	for _, value := range allEvents {
		value, _ := json.Marshal(value)
		fmt.Fprintf(w, string(value))
	}
}

func readHandler(w http.ResponseWriter, r *http.Request) {
	ID, err := r.URL.Query()["ID"]
	if !err {
		errorInput()
		return
	}
	for _, value := range allEvents {
		if value.ID == ID[0] {
			value, _ := json.Marshal(value)
			fmt.Fprintf(w, string(value))
		}
	}
}

func errorInput() {
	fmt.Printf("error input!")
}

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/create/", createHandler)
	http.HandleFunc("/read/", readHandler)
	http.HandleFunc("/readAll/", readAllHandler)
	log.Fatal(http.ListenAndServe(":8180", nil))
}
