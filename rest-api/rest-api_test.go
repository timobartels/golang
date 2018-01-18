package main

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigInit(t *testing.T) {

	ConfigInit()
	assert.Equal(t, port, ":8080")
}

func TestGetPeople(t *testing.T) {

	router := Routes()
	url := "/people"
	statusCode := 200
	people = append(people, Person{ID: "1", Firstname: "John", Lastname: "Doe"})

	req := httptest.NewRequest("GET", url, nil)
	req.Header.Set("content-type", "application/json")
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)

	assert.Equal(t, statusCode, res.Code)
	assert.JSONEq(t, string(res.Body.Bytes()), `[{"id":"1","firstname":"John","lastname":"Doe"}]`)
}

func TestCreatePerson(t *testing.T) {

	router := Routes()
	url := "/people/2"
	statusCode := 201

	data := map[string]string{"firstname": "Jane", "lastname": "Doe"}
	jsondata, _ := json.Marshal(data)

	req := httptest.NewRequest("POST", url, bytes.NewBuffer(jsondata))
	req.Header.Set("content-type", "application/json")
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)

	assert.Equal(t, statusCode, res.Code)
	assert.JSONEq(t, string(res.Body.Bytes()), `[{"id":"1","firstname":"John","lastname":"Doe"},{"id":"2","firstname":"Jane","lastname":"Doe"}]`)
}

func TestGetPerson(t *testing.T) {

	router := Routes()
	url := "/people/1"
	statusCode := 200

	req := httptest.NewRequest("GET", url, nil)
	req.Header.Set("content-type", "application/json")
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)

	assert.Equal(t, statusCode, res.Code)
	assert.JSONEq(t, string(res.Body.Bytes()), `{"id":"1","firstname":"John","lastname":"Doe"}`)
}

func TestDeletePerson(t *testing.T) {

	router := Routes()
	url := "/people/2"
	statusCode := 201

	req := httptest.NewRequest("DELETE", url, nil)
	req.Header.Set("content-type", "application/json")
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)

	assert.Equal(t, statusCode, res.Code)
	assert.NotContains(t, string(res.Body.Bytes()), `[{"id":"2","firstname":"Jane","lastname":"Doe"}]`)
}
