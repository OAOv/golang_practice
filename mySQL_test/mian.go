package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type Event struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Age  string `json:"age"`
}

var db *sql.DB
var err error

func getEvents(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var events []Event

	result, err := db.Query("SELECT * FROM user")
	if err != nil {
		panic(err.Error())
	}

	defer result.Close()

	for result.Next() {
		var event Event
		err := result.Scan(&event.ID, &event.Name, &event.Age)
		if err != nil {
			panic(err.Error())
		}
		events = append(events, event)
	}

	json.NewEncoder(w).Encode(events)
}

func createEvent(w http.ResponseWriter, r *http.Request) {
	stmt, err := db.Prepare("INSERT INTO user (name, age) VALUES (?, ?)")
	if err != nil {
		panic(err.Error())
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}

	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)
	name := keyVal["name"]
	age := keyVal["age"]

	_, err = stmt.Exec(name, age)
	if err != nil {
		panic(err.Error())
	}

	fmt.Fprintf(w, "New event was created.")
}

func getEvent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	result, err := db.Query("SELECT * FROM user WHERE id = ?", params["id"])
	if err != nil {
		panic(err.Error())
	}

	defer result.Close()

	var event Event
	for result.Next() {
		err := result.Scan(&event.ID, &event.Name, &event.Age)
		if err != nil {
			panic(err.Error())
		}
	}

	json.NewEncoder(w).Encode(event)
}

func updateEvent(w http.ResponseWriter, r *http.Request) {
}

func deleteEvent(w http.ResponseWriter, r *http.Request) {
}

func main() {
	db, err = sql.Open("mysql", "root:0000@tcp(127.0.0.1:3306)/test")

	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	router := mux.NewRouter()
	router.HandleFunc("/events", getEvents).Methods("GET")
	router.HandleFunc("/events", createEvent).Methods("POST")
	router.HandleFunc("/events/{id}", getEvent).Methods("GET")
	router.HandleFunc("/events/{id}", updateEvent).Methods("PUT")
	router.HandleFunc("/events/{id}", deleteEvent).Methods("DELETE")

	http.ListenAndServe(":8000", router)
}
