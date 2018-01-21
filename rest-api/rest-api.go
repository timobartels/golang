package main

import (
	"encoding/json"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
)

// Person struct to store the data for the API
type Person struct {
	ID        string `json:"id,omitempty"`
	Firstname string `json:"firstname,omitempty"`
	Lastname  string `json:"lastname,omitempty"`
}

var people []Person
var port string

// Defining the custom Prometheus metrics
// * http_requests_total_rest_api
// * http_request_duration
var (
	http_requests_total_rest_api = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total_rest_api",
			Help: "Number of http requests for rest-api.",
		},
		[]string{"host"},
	)
	http_request_duration_seconds = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "http_request_duration_seconds",
			Help: "Duration of the last request to rest-api.",
		},
		[]string{"host"},
	)
)

func main() {

	log.Info("HTTP server started and listening on port: ", port)

	// start our server
	log.Fatal(http.ListenAndServe(port, Routes()))

}

func init() {
	prometheus.MustRegister(http_request_duration_seconds)
	prometheus.MustRegister(http_requests_total_rest_api)
	ConfigInit()
}

// ConfigInit sets up configuration management using viper
func ConfigInit() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}
	port = viper.GetString("port")
	port = ":" + port
}

// Routes sets up the routes for our API
func Routes() http.Handler {
	router := mux.NewRouter()
	router.HandleFunc("/people", GetPeople).Methods("GET")
	router.HandleFunc("/people/{id}", GetPerson).Methods("GET")
	router.HandleFunc("/people/{id}", CreatePerson).Methods("POST")
	router.HandleFunc("/people/{id}", DeletePerson).Methods("DELETE")
	router.Path("/metrics").Handler(prometheus.Handler())
	return router
}

// GetPeople will output all entries in the people slice
func GetPeople(w http.ResponseWriter, r *http.Request) {
	requestStart := time.Now()
	http_requests_total_rest_api.With(prometheus.Labels{"host": "samplehost"}).Inc()
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(people)
	requestDuration := time.Since(requestStart).Seconds()
	http_request_duration_seconds.With(prometheus.Labels{"host": "samplehost"}).Set(float64(requestDuration))
}

// GetPerson will output only a specific entry
func GetPerson(w http.ResponseWriter, r *http.Request) {
	http_requests_total_rest_api.With(prometheus.Labels{"host": "samplehost"}).Inc()
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	params := mux.Vars(r)
	for _, item := range people {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			w.WriteHeader(http.StatusOK)
			return
		}
	}
}

// CreatePerson will create a new entry in the people slice
func CreatePerson(w http.ResponseWriter, r *http.Request) {
	http_requests_total_rest_api.With(prometheus.Labels{"host": "samplehost"}).Inc()
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	params := mux.Vars(r)
	var person Person
	_ = json.NewDecoder(r.Body).Decode(&person)
	person.ID = params["id"]
	people = append(people, person)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(people)
}

// DeletePerson will reshuffle the people slice to overwrite an entry
func DeletePerson(w http.ResponseWriter, r *http.Request) {
	http_requests_total_rest_api.With(prometheus.Labels{"host": "samplehost"}).Inc()
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	params := mux.Vars(r)
	for index, item := range people {
		if item.ID == params["id"] {
			people = append(people[:index], people[index+1:]...)
			break
		}
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(people)
}
