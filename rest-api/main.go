package main

import (
	log "github.com/sirupsen/logrus"

	"github.com/timobartels/golang/rest-api/model"
	"github.com/timobartels/golang/rest-api/server"
)

func main() {
	people := model.NewStore()

	restApp := server.NewRestApp(people)
	log.Fatal(restApp.Server().ListenAndServe())
}
