# Simple REST API in Go (tested with Go 1.9.2)

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

### Start the rest-server

* If you have a working Golang installation, you can simply use following command:  
```go run rest-api.go &```  

* You can compile the code via a Dockerfile and create/run a docker container of the rest-api server (requires Docker version that supports multi-stage builds, ie. Docker CE 17*):  
```docker build -t local/rest-server:1.0 . ```  
```docker run -d -p 8080:8080 local/rest-server:1.0```  



### Testing the API

Start the REST API in one terminal (or send it to background):
```go run rest-api.go```

Open another terminal window. You can test the REST API with the **cURL** command:

* Create a new record number 3:  
```curl -X POST -H "Content-Type:application/json" -d '{"firstname":"John", "lastname":"Doe"}' http://localhost:8080/people/1```

* Get a specific record by ID:  
```curl -X GET http://localhost:8080/people/2```

* Get all people defined:  
```curl -X GET http://localhost:8080/people```

* Delete record number 1:  
```curl -X DELETE -H "Content-Type:application/json" http://localhost:8080/people/1```
