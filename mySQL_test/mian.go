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

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Age  string `json:"age"`
}

var db *sql.DB
var err error

func getUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var users []User

	result, err := db.Query("SELECT * FROM user")
	if err != nil {
		panic(err.Error())
	}

	defer result.Close()

	for result.Next() {
		var user User
		err := result.Scan(&user.ID, &user.Name, &user.Age)
		if err != nil {
			panic(err.Error())
		}
		users = append(users, user)
	}

	json.NewEncoder(w).Encode(users)
}

func createUser(w http.ResponseWriter, r *http.Request) {
	stmt, err := db.Prepare("INSERT INTO user (name, age) VALUES (?, ?)")
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()

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

	fmt.Fprintf(w, "New User was created.")
}

func getUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	result, err := db.Query("SELECT * FROM user WHERE id = ?", params["id"])
	if err != nil {
		panic(err.Error())
	}

	defer result.Close()

	var user User
	for result.Next() {
		err := result.Scan(&user.ID, &user.Name, &user.Age)
		if err != nil {
			panic(err.Error())
		}
	}

	json.NewEncoder(w).Encode(user)
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	stmt, err := db.Prepare("UPDATE user SET name  = ?, age = ? WHERE id = ?")
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}

	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)
	newName, existsName := keyVal["name"]
	newAge, existsAge := keyVal["age"]

	if !existsName {
		stmt, err = db.Prepare("UPDATE user SET age = ? WHERE id = ?")
		_, err = stmt.Exec(newAge, params["id"])
		if err != nil {
			panic(err.Error())
		}
	} else if !existsAge {
		stmt, err = db.Prepare("UPDATE user SET name = ? WHERE id = ?")

		_, err = stmt.Exec(newName, params["id"])
		if err != nil {
			panic(err.Error())
		}
	} else {
		_, err = stmt.Exec(newName, newAge, params["id"])
		if err != nil {
			panic(err.Error())
		}
	}
	fmt.Fprintf(w, "User with ID = %s was updated", params["id"])
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	stmt, err := db.Prepare("DELETE FROM user WHERE id = ?")
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()

	_, err = stmt.Exec(params["id"])
	if err != nil {
		panic(err.Error())
	}

	fmt.Fprintf(w, "User with ID = %s was deleted", params["id"])
}

/*
func httpHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("./view/layout.html"))

	if r.Method != http.MethodPost {
		tmpl.Execute(w, struct {
			Title   string
			Input   []string
			Operate []string
		}{
			Title:   "users",
			Input:   []string{"ID", "Name", "Age"},
			Operate: []string{"create", "readAll", "readOne", "update", "delete"},
		})
	} else {
		log.Println("post something")
		user := User{
			ID:   r.FormValue("ID"),
			Name: r.FormValue("Name"),
			Age:  r.FormValue("Age"),
		}
		log.Println("ID: " + user.ID)

		tmpl.Execute(w, nil)
	}
}
*/
func main() {
	db, err = sql.Open("mysql", "root:0000@tcp(127.0.0.1:3306)/test")

	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	routerSQL := mux.NewRouter()
	routerSQL.HandleFunc("/users", getUsers).Methods("GET")
	routerSQL.HandleFunc("/users", createUser).Methods("POST")
	routerSQL.HandleFunc("/users/{id}", getUser).Methods("GET")
	routerSQL.HandleFunc("/users/{id}", updateUser).Methods("PATCH")
	routerSQL.HandleFunc("/users/{id}", deleteUser).Methods("DELETE")

	go http.ListenAndServe(":8000", routerSQL)

	/*
		http.HandleFunc("/", httpHandler)
		http.ListenAndServe(":8010", nil)
	*/
}
