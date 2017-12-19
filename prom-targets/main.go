package main

import (
	"fmt"
	"os"

	"github.com/timobartels/golang/prom-targets/prom"
)

func main() {

	promUrl := "http://192.168.56.10:9090"

	targets, err := prom.GetTargets(promUrl)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(targets)
}
