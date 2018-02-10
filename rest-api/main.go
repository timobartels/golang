package main

import (
	"github.com/timobartels/golang/rest-api/server"
)

func main() {

	// Creating application instance
	a := server.App{}

	// Initialize the REST server
	a.Initialize()

	// Start our server
	a.Run()
}
