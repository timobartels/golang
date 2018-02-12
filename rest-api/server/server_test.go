package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createDummyEntry() *RestApp {

	people := NewStore()
	app := NewRestApp(people)

	data := map[string]string{"firstname": "Jane", "lastname": "Doe"}
	jsondata, _ := json.Marshal(data)
	req, _ := http.NewRequest("POST", "/people/1", bytes.NewBuffer(jsondata))
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(app.CreatePerson)
	handler.ServeHTTP(rr, req)
	return &app
}

func TestConfigInitJSONWarn(t *testing.T) {

	var configFile string
	configFile = "configjson"
	ConfigInit(configFile)
}

func TestConfigInitTextDebug(t *testing.T) {

	var configFile string
	configFile = "configdebug"
	ConfigInit(configFile)
}

func TestConfigInitError(t *testing.T) {

	var configFile string
	configFile = "missing"
	ConfigInit(configFile)
}

func TestGetPeople(t *testing.T) {

	people := NewStore()
	app := NewRestApp(people)

	req, err := http.NewRequest("GET", "/people", nil)
	if err != nil {
		t.FailNow()
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(app.GetPeople)
	handler.ServeHTTP(rr, req)
	assert.Equal(t, 200, rr.Code)
}

func TestCreatePerson(t *testing.T) {

	people := NewStore()
	app := NewRestApp(people)

	data := map[string]string{"firstname": "Jane", "lastname": "Doe"}
	jsondata, _ := json.Marshal(data)

	req, err := http.NewRequest("POST", "/people/1", bytes.NewBuffer(jsondata))
	if err != nil {
		t.FailNow()
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(app.CreatePerson)
	handler.ServeHTTP(rr, req)
	assert.Equal(t, 201, rr.Code)
	assert.JSONEq(t, `{"Store":[{"firstname":"Jane","lastname":"Doe"}]}`, string(rr.Body.Bytes()))
}

func TestGetPerson(t *testing.T) {

	app := createDummyEntry()

	req, err := http.NewRequest("GET", "/people/1", nil)
	if err != nil {
		t.FailNow()
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(app.GetPerson)
	handler.ServeHTTP(rr, req)
	assert.Equal(t, 200, rr.Code)
	assert.JSONEq(t, `{"firstname":"Jane","lastname":"Doe"}`, string(rr.Body.Bytes()))
}

func TestDeletePerson(t *testing.T) {

	app := createDummyEntry()

	req, err := http.NewRequest("DELETE", "/people/1", nil)
	if err != nil {
		t.FailNow()
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(app.DeletePerson)
	handler.ServeHTTP(rr, req)
	assert.Equal(t, 201, rr.Code)
	assert.JSONEq(t, `{"Store":[]}`, string(rr.Body.Bytes()))
}

func TestServerCreated(t *testing.T) {

	people := NewStore()
	app := NewRestApp(people)

	assert.NotNil(t, app.Server())
}
