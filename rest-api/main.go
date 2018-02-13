package main

import (
	log "github.com/sirupsen/logrus"

	"github.com/timobartels/golang/rest-api/server"
)

var config string

func main() {
	backend := server.NewStore()
	config = "config"
	server.ConfigInit(config)

	restApp := server.NewRestApp(backend)
	log.Fatal(restApp.Server().ListenAndServe())
}
