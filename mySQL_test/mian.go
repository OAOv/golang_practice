package main

import (
	"database/sql"
	"encoding/json"
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
	w.Header().Set("Content-Type", "application/json")

}

func getEvent(w http.ResponseWriter, r *http.Request) {
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
