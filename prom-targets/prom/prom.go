package prom

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"
)

type PromTargets struct {
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

func GetTargets(promUrl string) (PromTargets, error) {

	url := promUrl + "/api/v1/targets"
	data := PromTargets{}
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return data, errors.New("Unable to connect to the Prometheus server")
	}
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	response, err := client.Do(request)
	defer response.Body.Close()

	if err != nil {
		return data, errors.New("Unable to access the Prometheus targets endpoint")
	}

	jsonresp, _ := ioutil.ReadAll(response.Body)

	jsonerr := json.Unmarshal(jsonresp, &data)
	if jsonerr != nil {
		return data, jsonerr
	}
	return data, nil
}
