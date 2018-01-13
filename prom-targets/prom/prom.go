// Package prom provides a function to read discovered targets from Prometheus server
package prom

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

func LogInit(logFormat string) {
	// Log as JSON instead of the default ASCII formatter.
	if logFormat == "json" {
		log.SetFormatter(&log.JSONFormatter{})
	}

	// Output to stdout instead of the default stderr
	log.SetOutput(os.Stdout)
	// Only log the warning severity or above.
	log.SetLevel(log.DebugLevel)
}

// Targets custom type keeps the JSON decoded from Prometheus /api/v1/targets endpoint
type Targets struct {
	Status string `json:"status"`
	Data   struct {
		ActiveTargets []struct {
			DiscoveredLabels struct {
				Address       string `json:"__address__"`
				MetricsPath   string `json:"__metrics_path__"`
				MarathonApp   string `json:"__meta_marathon_app"`
				MarathonImage string `json:"__meta_marathon_image"`
				Scheme        string `json:"__scheme__"`
				Job           string `json:"job"`
			} `json:"discoveredLabels"`
			Labels struct {
				AppID    string `json:"app_id"`
				Instance string `json:"instance"`
				Job      string `json:"job"`
			} `json:"labels"`
			ScrapeURL  string    `json:"scrapeUrl"`
			LastError  string    `json:"lastError"`
			LastScrape time.Time `json:"lastScrape"`
			Health     string    `json:"health"`
		} `json:"activeTargets"`
	} `json:"data"`
}

// GetTargets reads all discovered targets from provided Prometheus URL /api/v1/targets endpoint.
// It will return the custom type PromTargets to the calling function for further processing.
func GetTargets(promURL string) (Targets, error) {

	url := promURL + "/api/v1/targets"
	data := Targets{}
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return data, err
	}
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	response, err := client.Do(request)
	if err != nil {
		return data, errors.New("Unable to access the Prometheus targets endpoint")
	}
	defer response.Body.Close()

	jsonresp, _ := ioutil.ReadAll(response.Body)

	jsonerr := json.Unmarshal(jsonresp, &data)
	if jsonerr != nil {
		return data, jsonerr
	}
	for _, v := range data.Data.ActiveTargets {
		log.WithFields(log.Fields{
			"instance":      v.Labels.Instance,
			"marathon-app":  v.DiscoveredLabels.MarathonApp,
			"target-health": v.Health,
			"promURL":       promURL,
			"status":        "up",
		}).Info("Output from Prometheus response")
	}
	return data, nil
}
