package main

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/timobartels/golang/prom-targets/prom"
)

func main() {

	promURL := "http://192.168.56.10:9090"

	prom.LogInit()

	targets, err := prom.GetTargets(promURL)
	if err != nil {
		//		fmt.Println(err)
		log.WithFields(log.Fields{
			"promURL": promURL,
			"status":  "down",
		}).Fatal(err)
		os.Exit(1)
	}

	log.WithFields(log.Fields{
		"promURL": promURL,
		"status":  "up",
	}).Info("Successful response from Prometheus API")

	fmt.Println(targets)
	/*		for _, v := range targets.Data.ActiveTargets {
				if v.Labels.Job == "DCOS Services Monitoring" {
					fmt.Println("TaskID:", v.DiscoveredLabels.MarathonTask, "Health:", v.Health)
				}
			}
	*/
}
