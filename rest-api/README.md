# Simple REST API in Go

This is a simple (demo) REST server with very minimal functionality:
* get all person records defined so far (GET)
* get a specific person record (GET)
* create a new person record (POST)
* delete a specific person record (DELETE)

This go program utilizes following modules:
* gorilla/mux (for request routing)
* logrus (for logging)
* viper (for configuration management)
* prometheus (for exposing a Prometheus metrics monitoring endpoint)

The program has been developed utilizing test-driven-development (TDD)


### Testing the API

Start the REST API in one terminal (or send it to background):
```go run rest-api.go```

Open another terminal window. You can test the REST API with the **cURL** command:

* Create a new record number 3:  
```curl -X POST -H "Content-Type:application/json" -d '{"firstname":"John", "lastname":"Doe"}' http://localhost:8080/people/1```

* Get all people defined:  
```curl -X GET http://localhost:8080/people```

* Delete record number 1:  
```curl -X DELETE -H "Content-Type:application/json" http://localhost:8080/people/1```
