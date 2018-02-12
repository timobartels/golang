package main

import (
	log "github.com/sirupsen/logrus"

	"github.com/timobartels/golang/rest-api/server"
)

var config string

func main() {
	people := server.NewStore()
	config = "config"
	server.ConfigInit(config)

	restApp := server.NewRestApp(people)
	log.Fatal(restApp.Server().ListenAndServe())
}
