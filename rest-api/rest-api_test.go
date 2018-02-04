package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

var app App

func TestGetPeople(t *testing.T) {

	app = App{}
	app.Initialize()
	req, err := http.NewRequest("GET", "/people", nil)
	if err != nil {
		t.FailNow()
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(app.GetPeople)
	handler.ServeHTTP(rr, req)
	assert.Equal(t, rr.Code, 200, rr.Code)
}

func TestCreatePerson(t *testing.T) {

	data := map[string]string{"firstname": "Jane", "lastname": "Doe"}
	jsondata, _ := json.Marshal(data)

	req, err := http.NewRequest("POST", "/people/1", bytes.NewBuffer(jsondata))
	if err != nil {
		t.FailNow()
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(app.CreatePerson)
	handler.ServeHTTP(rr, req)
	assert.Equal(t, rr.Code, 201, rr.Code)
	assert.JSONEq(t, string(rr.Body.Bytes()), `[{"firstname":"Jane","lastname":"Doe"}]`)
}

func TestGetPerson(t *testing.T) {

	req, err := http.NewRequest("GET", "/people/1", nil)
	if err != nil {
		t.FailNow()
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(app.GetPerson)
	handler.ServeHTTP(rr, req)
	assert.Equal(t, rr.Code, 200, rr.Code)
}

func TestDeletePerson(t *testing.T) {

	req, err := http.NewRequest("DELETE", "/people/1", nil)
	if err != nil {
		t.FailNow()
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(app.DeletePerson)
	handler.ServeHTTP(rr, req)
	assert.Equal(t, rr.Code, 201, rr.Code)
	assert.NotContains(t, string(rr.Body.Bytes()), `[{"id":"1","firstname":"Jane","lastname":"Doe"}]`)
}
