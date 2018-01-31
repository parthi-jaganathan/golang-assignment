# golang-assignment project

My first Go Project.

golang-assignment project provides a simple webserver that provides endpoints to hash password using SHA512 and get statictics about the hash endpoints

List of available service endpoints and their functionality

* Generate Sequence ID for a given password hash request
    - POST `/hash`
    - Payload format password={password}
    - ```curl -X POST -d "password=angryMonkey" http://localhost:8080/hash```
    - The API returns an sequence ID
    - The sequence ID can be used to retrieve the password hash using the below endpoint
    - Uses in memory map to track the sequence ID and password Hash
* Retrieve password hash using Sequece ID from above endpoint
    - GET `/hash/{id}`
    - Use the sequence ID from the previous POST call to get the hash
    - ```curl -X GET http://localhost:8080/hash/2```
* Get statictics of above endpoints
    - GET `/stats`
    - ```curl -X GET http://localhost:8080/stats```
    - Sample Response: { "total":4,"average":1250 }
        - `total` -> total number of successful requests made to /hash and /hash/{id} endpoints
        - `average` -> captures the average response time for the endpoints /hash and /hash/{id} in milliseconds (successful requests)
    - Uses in memory map to track the statictics
* Graceful shutdown of the server
    - GET `/shutdown`
    - ```curl -X GET http://localhost:8080/shutdown```
    - Return OK and waits for the active requests to finish processing before gracefully shutting down the server

# To get the project

This will download the project into your $HOME/go/src directory.

```
go get github.com/pjaganathan/golang-assignment
```

# To Run

After getting the project you can now run the below commands in sequence to build and run the server

```
cd $HOME/go/src/github.com/pjaganathan/golang-assignment
go build
go install
$HOME/go/bin/golang-assignment

```

OR you can run directly from the golang-assignment checkout directory

```
cd $HOME/go/src/github.com/pjaganathan/golang-assignment
go build
./golang-assignment
```

The server will run on port 8080 which is configured in main.go

