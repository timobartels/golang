package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

const succeed = "\u2713"
const failed = "\u2717"

func getPeopleRoute() http.Handler {

	r := mux.NewRouter()
	r.HandleFunc("/people", GetPeople).Methods("GET")
	return r
}

func TestGetPeople(t *testing.T) {

	url := "/people"
	statusCode := 200
	people = append(people, Person{ID: "1", Firstname: "John", Lastname: "Doe"})

	t.Log("Given the need to test the GetPeople endpoint.")
	{
		req := httptest.NewRequest("GET", url, nil)
		req.Header.Set("content-type", "application/json")
		res := httptest.NewRecorder()
		getPeopleRoute().ServeHTTP(res, req)

		fmt.Println(res)

		t.Logf("\tTest 0:\tWhen checking %q for status code %d", url, statusCode)
		{
			if res.Code != 200 {
				t.Fatalf("\t%s\tShould receive a status code of %d for the response. Received[%d].", failed, statusCode)
			}
			t.Logf("\t%s\tShould receive a status code of %d for the response.", succeed, statusCode)

			var u Person

			if err := json.NewDecoder(res.Body).Decode(&u); err != nil {
				t.Fatalf("\t%s\tShould be able to decode the response. Error was: %q", failed, err)
			}
			t.Logf("\t%s\tShould be able to decode the response.", succeed)

		}
	}
}

func createPersonRoute() http.Handler {

	r := mux.NewRouter()
	r.HandleFunc("/people/{id}", CreatePerson).Methods("POST")
	return r
}

func TestCreatePerson(t *testing.T) {

	url := "/people/2"
	statusCode := 200

	data := map[string]string{"firstname": "Jane", "lastname": "Doe"}
	jsondata, _ := json.Marshal(data)

	t.Log("Given the need to test the CreatePerson endpoint.")
	{
		req := httptest.NewRequest("POST", url, bytes.NewBuffer(jsondata))
		req.Header.Set("content-type", "application/json")
		res := httptest.NewRecorder()
		createPersonRoute().ServeHTTP(res, req)

		t.Logf("\tTest 0:\tWhen checking %q for status code %d", url, statusCode)
		{
			if res.Code != 200 {
				t.Fatalf("\t%s\tShould receive a status code of %d for the response. Received[%d].", failed, statusCode)
			}
			t.Logf("\t%s\tShould receive a status code of %d for the response.", succeed, statusCode)

			var u Person

			fmt.Println(res)

			if err := json.NewDecoder(res.Body).Decode(&u); err != nil {
				t.Fatalf("\t%s\tShould be able to decode the response. Error was: %q", failed, err)
			}
			t.Logf("\t%s\tShould be able to decode the response.", succeed)

			if u.Firstname == "Mike" {
				t.Logf("\t%s\tShould have \"Mike\" for Firstname in the response.", succeed)
			} else {
				t.Fatalf("\t%s\tShould have \"Mike\" for Firstname in the response : %q", failed, u.Firstname)
			}
		}
	}
}
