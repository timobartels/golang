package server

import (
	"encoding/json"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/timobartels/golang/rest-api/model"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
)

// RestApp is holding the rest logic and handlers
type RestApp struct {
	people *model.People
}

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

func NewRestApp(people *model.People) RestApp {
	return RestApp{people}
}

func (app *RestApp) Server() *http.Server {
	router := mux.NewRouter()
	router.HandleFunc("/people", app.GetPeople).Methods("GET")
	router.HandleFunc("/people/{id}", app.GetPerson).Methods("GET")
	router.HandleFunc("/people/{id}", app.CreatePerson).Methods("POST")
	router.HandleFunc("/people/{id}", app.DeletePerson).Methods("DELETE")
	router.Path("/metrics").Handler(prometheus.Handler())

	server := &http.Server{
		Handler:      router,
		Addr:         "127.0.0.1" + port,
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  30 * time.Second,
	}
	return server
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

// GetPeople will output all entries in the people slice
func (app *RestApp) GetPeople(w http.ResponseWriter, r *http.Request) {
	log.WithFields(log.Fields{"method": r.Method, "endpoint": r.URL.Path}).Info("Reveived new request.")
	requestStart := time.Now()
	http_requests_total_rest_api.With(prometheus.Labels{"path": r.URL.Path, "method": r.Method}).Inc()
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(app.people)
	requestDuration := time.Since(requestStart).Seconds()
	http_request_duration_seconds.With(prometheus.Labels{"path": r.URL.Path, "method": r.Method}).Set(float64(requestDuration))
}

// GetPerson will output only a specific entry
func (app *RestApp) GetPerson(w http.ResponseWriter, r *http.Request) {
	log.WithFields(log.Fields{"method": r.Method, "endpoint": r.URL.Path}).Info("Reveived new request.")
	requestStart := time.Now()
	http_requests_total_rest_api.With(prometheus.Labels{"path": r.URL.Path, "method": r.Method}).Inc()
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	params := mux.Vars(r)
	for _, item := range app.people.Store {
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
func (app *RestApp) CreatePerson(w http.ResponseWriter, r *http.Request) {
	log.WithFields(log.Fields{"method": r.Method, "endpoint": r.URL.Path}).Info("Reveived new request.")
	requestStart := time.Now()
	http_requests_total_rest_api.With(prometheus.Labels{"path": r.URL.Path, "method": r.Method}).Inc()
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	params := mux.Vars(r)
	var person model.Person
	_ = json.NewDecoder(r.Body).Decode(&person)
	person.ID = params["id"]
	app.people.Store = append(app.people.Store, person)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(app.people)
	requestDuration := time.Since(requestStart).Seconds()
	http_request_duration_seconds.With(prometheus.Labels{"path": r.URL.Path, "method": r.Method}).Set(float64(requestDuration))
}

// DeletePerson will reshuffle the people slice to overwrite an entry
func (app *RestApp) DeletePerson(w http.ResponseWriter, r *http.Request) {
	log.WithFields(log.Fields{"method": r.Method, "endpoint": r.URL.Path}).Info("Reveived new request.")
	requestStart := time.Now()
	http_requests_total_rest_api.With(prometheus.Labels{"path": r.URL.Path, "method": r.Method}).Inc()
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	params := mux.Vars(r)
	for index, item := range app.people.Store {
		if item.ID == params["id"] {
			app.people.Store = append(app.people.Store[:index], app.people.Store[index+1:]...)
			break
		}
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(app.people)
	requestDuration := time.Since(requestStart).Seconds()
	http_request_duration_seconds.With(prometheus.Labels{"path": r.URL.Path, "method": r.Method}).Set(float64(requestDuration))
}
