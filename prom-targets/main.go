package main

import (
	"fmt"
	"os"

	"github.com/timobartels/golang/prom-targets/prom"
)

func main() {

	promURL := "http://192.168.56.10:9090"

	targets, err := prom.GetTargets(promURL)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(targets)
	for _, v := range targets.Data.ActiveTargets {
		if v.Labels.Job == "DCOS Services Monitoring" {
			fmt.Println("TaskID:", v.DiscoveredLabels.MarathonTask, "Health:", v.Health)
		}
	}
}
