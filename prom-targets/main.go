package main

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/timobartels/golang/prom-targets/prom"
)

func main() {

	// setup config management
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()

	if err != nil {
		log.Info("No config file found, using default values")
	}
	// set defaults
	viper.SetDefault("promURL", "http://192.168.56.10:9090")

	prom.LogInit()

	promURL := viper.GetString("promURL")

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
