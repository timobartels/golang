package main

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/gorilla/mux"
)

// Person struct to store the data for the API
type Person struct {
	ID        string `json:"id,omitempty"`
	Firstname string `json:"firstname,omitempty"`
	Lastname  string `json:"lastname,omitempty"`
}

var people []Person

func main() {

	// setup config management
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("No config file found, terminating program!")
	}
	// get details from config file
	port := viper.GetString("port")
	port = ":" + port

	// fill out people slice with some dummy data
	people = append(people, Person{ID: "1", Firstname: "Mike", Lastname: "Brewers"})
	people = append(people, Person{ID: "2", Firstname: "Jennifer", Lastname: "Myers"})

	// initialize routes
	r := Routes()

	log.Info("HTTP server started and listening on port: ", port)

	// start our server
	log.Fatal(http.ListenAndServe(port, r))

}

// Routes sets up the routes for our API
//func Routes() http.Handler {
func Routes() http.Handler {
	router := mux.NewRouter()
	router.HandleFunc("/people", GetPeople).Methods("GET")
	router.HandleFunc("/people/{id}", GetPerson).Methods("GET")
	router.HandleFunc("/people/{id}", CreatePerson).Methods("POST")
	router.HandleFunc("/people/{id}", DeletePerson).Methods("DELETE")
	return router
}

// GetPeople will output all entries in the people slice
func GetPeople(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(people)
}

// GetPerson will output only a specific entry
func GetPerson(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for _, item := range people {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
}

// CreatePerson will create a new entry in the people slice
func CreatePerson(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var person Person
	_ = json.NewDecoder(r.Body).Decode(&person)
	person.ID = params["id"]
	people = append(people, person)
	json.NewEncoder(w).Encode(people)
}

// DeletePerson will reshuffle the people slice to overwrite an entry
func DeletePerson(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for index, item := range people {
		if item.ID == params["id"] {
			people = append(people[:index], people[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(people)
}
