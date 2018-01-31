# GoLang WebServer

My first Go Project.

golang-assignment project provides a simple webserver that provides endpoints to hash password using SHA512 and get statictics about the hash endpoints

List of available service endpoints and their functionality

* POST /hash
** Payload format password={password}
** Eg: curl -X POST -d "password=angryMonkey" http://localhost:8080/hash
** The API returns an sequence Identifier
** The sequence Identifier can be used to retrieve the password Hash using the below endpoint
** Uses in memory map to track the sequence ID and password Hash
* GET /hash/{id}
** Use the sequence ID from the previous POST call to get the hash
** Eg: curl -X GET http://localhost:8080/hash/2
* GET /stats
** Eg: curl -X GET http://localhost:8080/stats
** Sample Response 
{"total":4,"average":1250}
** total -> total number of successful requests made to /hash and /hash{id} endpoints
** average -> captures the average response time for the endpoints /hash and /hash/{id} (successful requests)
** Uses in memory map to track the statictics
* GET /shutdown
** Eg: curl -X GET http://localhost:8080/shutdown
** Return OK and waits for the active requests to finish processing before gracefully shutting down the server

# To Get this assignment

> go get github.com/pjaganathan/golang-assignment

# To Run

Once you have the project in your $HOME/go/src directory, you can then do the below to run it

> cd $HOME/go/src/github.com/pjaganathan/golang-assignment
> go build
> go install
> $HOME/go/bin/golang-assignment

The server will run on port 8080 which is configured in main.go

