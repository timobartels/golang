package prom

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

const succeed = "\u2713"
const failed = "\u2717"

func mockTargetsEndpoint() *httptest.Server {

	var payload = `{"status":"success","data":{"activeTargets":[{"discoveredLabels":{"__address__":"10.20.30.40:1234","__meta_marathon_app":"/test/mytest-dev","__metrics_path__":"/metrics","job":"DCOS Services Monitoring"},"labels":{"app_id":"/test/mytest-dev","instance":"10.20.30.40:1234","job":"DCOS Services Monitoring"},"scrapeUrl":"http://10.20.30.40:1234/metrics","lastError":"","lastScrape":"2017-12-18T06:56:47.0300433Z","health":"up"},{"discoveredLabels":{"__address__":"11.21.31.41:5678","__meta_marathon_app":"/test/mytest-qa","__metrics_path__":"/metrics","job":"DCOS Services Monitoring"},"labels":{"app_id":"/test/mytest-qa","instance":"11.21.31.41:5678","job":"DCOS Services Monitoring"},"scrapeUrl":"http://11.21.31.41:5678/metrics","lastError":"","lastScrape":"2017-12-18T06:56:47.0300433Z","health":"down"}]}}`

	f := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, payload)
	}
	return httptest.NewServer(http.HandlerFunc(f))
}

func TestGetTargets(t *testing.T) {

	server := mockTargetsEndpoint()
	defer server.Close()

	t.Log("Given the need to test GetTargets function")

	t.Logf("\tTest 1: When calling the GetTargets function.")

	targets, err := GetTargets(server.URL)

	fmt.Println(targets)

	if err != nil {
		t.Fatalf("\t%s\tShould be able to connect to Prometheus and unmarshal JSON body: %v ", failed, err)
	}
	t.Logf("\t%s\tShould be able to connect to Prometheus and unmarshal JSON body.", succeed)

	t.Logf("\tTest 2: When reading the JSON response from Prometheus target endpoint.")

	expectedHealth := "up"

	if targets.Data.ActiveTargets[0].Health == expectedHealth {
		t.Logf("\t%s\tShould get the expected content: '%s'. Got: '%s' ", succeed, expectedHealth, targets.Data.ActiveTargets[0].Health)
	} else {
		t.Fatalf("\t%s\tShould get the expected content: '%s'. Got: '%v' ", failed, expectedHealth, err)
	}
}
