package httphandler

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/pjaganathan/golang-assignment/handlers/password"
	"github.com/pjaganathan/golang-assignment/handlers/stats"
	"github.com/pjaganathan/golang-assignment/utils"
)

type errorResponse struct {
	statusCode   int
	errorMessage string
}

const (
	passwordPrefix           = "password="
	invalidRequestPayloadMsg = "Invalid request payload"
	errorParsingPayloadMsg   = "Error when reading the request body"
	contentTypeHeader        = "Content-Type"
	applicationJSON          = "application/json"
)

func setErrorResponse(w http.ResponseWriter, err errorResponse) {
	const defaultStatusCode = http.StatusInternalServerError
	if err.statusCode == 0 {
		err.statusCode = http.StatusInternalServerError
	}
	if err.errorMessage == "" {
		err.errorMessage = "Unknown"
	}
	http.Error(w, err.errorMessage, err.statusCode)
}

// setResponseCodeAndMsg
func setResponseCodeAndCustomMsg(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	w.Write([]byte(message))
}

// setResponseCodeAndMsg
func setResponseCodeAndDefaultMsg(w http.ResponseWriter, statusCode int) {
	setResponseCodeAndCustomMsg(w, statusCode, http.StatusText(statusCode))
}

// RootHandler return 404 for any request to root of the application
func RootHandler() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		setResponseCodeAndDefaultMsg(w, http.StatusNotFound)
	}
	return http.HandlerFunc(fn)
}

// GetPasswordHashBySequenceID returns the password hash based on the request sequence ID
func GetPasswordHashBySequenceID(path string) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {

		// ensure the method is GET
		if r.Method != http.MethodGet {
			setResponseCodeAndDefaultMsg(w, http.StatusMethodNotAllowed)
			return
		}

		defer statsHandler.TrackHashPasswordAPIMetrics(time.Now()) // tracks the total requests and response time of the API
		idStr := strings.TrimPrefix(r.URL.Path, path)
		id, err := strconv.Atoi(idStr)
		if err != nil { // if not an integer, return 404
			setResponseCodeAndDefaultMsg(w, http.StatusNotFound)
			return
		}

		passHash, err := passwordHandler.GetPasswordHashBySequenceID(id)
		if err != nil {
			setResponseCodeAndDefaultMsg(w, http.StatusNotFound)
			return
		}

		// OK, return the hashed password
		time.Sleep(5 * time.Second) // sleeping for 5 seconds
		setResponseCodeAndCustomMsg(w, http.StatusOK, passHash)
	}
	return http.HandlerFunc(fn)
}

// GenerateRequestSequenceID returns a sequence ID that can be later used to retrieve the password hash
func GenerateRequestSequenceID() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {

		var err errorResponse
		defer func() {
			if r := recover(); r != nil {
				switch t := r.(type) {
				case string:
					err = errorResponse{errorMessage: t}
				case error:
					err = errorResponse{errorMessage: t.Error()}
				case errorResponse:
					err = t
				default:
					err = errorResponse{}
				}
				setErrorResponse(w, err)
			}
		}()

		if r.Method != http.MethodPost {
			log.Printf("Requested method %s is not allowed", r.Method)
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		defer statsHandler.TrackHashPasswordAPIMetrics(time.Now()) // tracks the total requests and response time of the API

		reqBody, readErr := ioutil.ReadAll(r.Body)
		if readErr != nil {
			setResponseCodeAndCustomMsg(w, http.StatusInternalServerError, readErr.Error())
			return
		}

		payload := string(reqBody) // validate if the request body has required format
		if strings.Index(payload, passwordPrefix) != 0 {
			setResponseCodeAndCustomMsg(w, http.StatusBadRequest, invalidRequestPayloadMsg)
			return
		}

		password := strings.TrimPrefix(payload, passwordPrefix) // split everything after prefix as one string
		if len(password) < 1 {
			setResponseCodeAndCustomMsg(w, http.StatusBadRequest, invalidRequestPayloadMsg)
			return
		}

		passHash, hashErr := passwordutil.GeneratePasswordHash(password) // generate the hash
		if hashErr != nil {
			setResponseCodeAndCustomMsg(w, http.StatusInternalServerError, hashErr.Error())
			return
		}

		sequenceID := strconv.Itoa(passwordHandler.GenerateSequenceID(passHash)) // generate sequenceID if everything went ok
		log.Printf("sending request sequenceID %v \n", sequenceID)
		setResponseCodeAndCustomMsg(w, http.StatusOK, sequenceID)
		w.Header().Set("location", r.RequestURI+"/"+sequenceID)

	}
	return http.HandlerFunc(fn)
}

// ShutdownGracefully accepts the request to graefully shutdown the server; Returns OK 201 immediately
func ShutdownGracefully(done chan bool) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {

		// ensure the method is GET
		if r.Method != http.MethodGet {
			setResponseCodeAndDefaultMsg(w, http.StatusMethodNotAllowed)
			return
		}

		setResponseCodeAndDefaultMsg(w, http.StatusAccepted)
		done <- true
	}
	return http.HandlerFunc(fn)
}

// GetStats returns the statistics of the /hash and /hash/{id} endpoint capturing total requests and average response times
func GetStats() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {

		// ensure the method is GET
		if r.Method != http.MethodGet {
			setResponseCodeAndDefaultMsg(w, http.StatusMethodNotAllowed)
			return
		}

		resp := statsHandler.GetStats()
		rs, err := json.Marshal(resp)

		if err != nil {
			log.Println("Stats handler error when parsing response data as json", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Println("Stats data ", string(rs))

		w.Header().Set(contentTypeHeader, applicationJSON)
		w.Write(rs)
	}
	return http.HandlerFunc(fn)
}
