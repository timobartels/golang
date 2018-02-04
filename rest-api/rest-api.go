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

// App is the application
type App struct {
	Router *mux.Router
}

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
		[]string{"path", "method"},
	)
	http_request_duration_seconds = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "http_request_duration_seconds",
			Help: "Duration of the last request to rest-api.",
		},
		[]string{"path", "method"},
	)
)

func main() {

	// Creating application instance
	a := App{}

	// Initialize the REST server
	a.Initialize()

	// Start our server
	a.Run(port)
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
	viper.SetDefault("port", "8080")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}
	port = viper.GetString("port")
	port = ":" + port
	loglevel := viper.GetString("loglevel")
	switch loglevel {
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "debug":
		log.SetLevel(log.DebugLevel)
	}
	logformat := viper.GetString("logformat")
	if logformat == "json" {
		log.SetFormatter(&log.JSONFormatter{})
		log.Info("Logging to stdout in JSON format.")
	}
}

// InitRoutes initializes the routes for the REST server
func (a *App) InitRoutes() {
	a.Router.HandleFunc("/people", a.GetPeople).Methods("GET")
	a.Router.HandleFunc("/people/{id}", a.GetPerson).Methods("GET")
	a.Router.HandleFunc("/people/{id}", a.CreatePerson).Methods("POST")
	a.Router.HandleFunc("/people/{id}", a.DeletePerson).Methods("DELETE")
	a.Router.Path("/metrics").Handler(prometheus.Handler())
	log.Info("Routes initialized.")
}

func (a *App) Run(port string) {
	log.Info("Starting HTTP server on port: ", port)
	log.Fatal(http.ListenAndServe(port, a.Router))
}

func (a *App) Initialize() {
	a.Router = mux.NewRouter()
	a.InitRoutes()
}

// GetPeople will output all entries in the people slice
func (a *App) GetPeople(w http.ResponseWriter, r *http.Request) {
	log.WithFields(log.Fields{"method": r.Method, "endpoint": r.URL.Path}).Info("Reveived new request.")
	requestStart := time.Now()
	http_requests_total_rest_api.With(prometheus.Labels{"path": r.URL.Path, "method": r.Method}).Inc()
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(people)
	requestDuration := time.Since(requestStart).Seconds()
	http_request_duration_seconds.With(prometheus.Labels{"path": r.URL.Path, "method": r.Method}).Set(float64(requestDuration))
}

// GetPerson will output only a specific entry
func (a *App) GetPerson(w http.ResponseWriter, r *http.Request) {
	log.WithFields(log.Fields{"method": r.Method, "endpoint": r.URL.Path}).Info("Reveived new request.")
	requestStart := time.Now()
	http_requests_total_rest_api.With(prometheus.Labels{"path": r.URL.Path, "method": r.Method}).Inc()
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	params := mux.Vars(r)
	for _, item := range people {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			w.WriteHeader(http.StatusOK)
			requestDuration := time.Since(requestStart).Seconds()
			http_request_duration_seconds.With(prometheus.Labels{"path": r.URL.Path, "method": r.Method}).Set(float64(requestDuration))
			return
		}
	}
}

// CreatePerson will create a new entry in the people slice
func (a *App) CreatePerson(w http.ResponseWriter, r *http.Request) {
	log.WithFields(log.Fields{"method": r.Method, "endpoint": r.URL.Path}).Info("Reveived new request.")
	requestStart := time.Now()
	http_requests_total_rest_api.With(prometheus.Labels{"path": r.URL.Path, "method": r.Method}).Inc()
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	params := mux.Vars(r)
	var person Person
	_ = json.NewDecoder(r.Body).Decode(&person)
	person.ID = params["id"]
	people = append(people, person)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(people)
	requestDuration := time.Since(requestStart).Seconds()
	http_request_duration_seconds.With(prometheus.Labels{"path": r.URL.Path, "method": r.Method}).Set(float64(requestDuration))
}

// DeletePerson will reshuffle the people slice to overwrite an entry
func (a *App) DeletePerson(w http.ResponseWriter, r *http.Request) {
	log.WithFields(log.Fields{"method": r.Method, "endpoint": r.URL.Path}).Info("Reveived new request.")
	requestStart := time.Now()
	http_requests_total_rest_api.With(prometheus.Labels{"path": r.URL.Path, "method": r.Method}).Inc()
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
	requestDuration := time.Since(requestStart).Seconds()
	http_request_duration_seconds.With(prometheus.Labels{"path": r.URL.Path, "method": r.Method}).Set(float64(requestDuration))
}
