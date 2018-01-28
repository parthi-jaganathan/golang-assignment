package httphandler

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"../utils"
)

type errorResponse struct {
	statusCode   int
	errorMessage string
}

const passwordPrefix = "password="

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

// RootHandler return 200 and blank response when root of the application is accessed
func RootHandler(path string) http.Handler {
	log.Printf("Initiating Root handler for path %s", path)
	handler := func(w http.ResponseWriter, r *http.Request) {
		setResponseCodeAndDefaultMsg(w, http.StatusOK)
	}
	return http.HandlerFunc(handler)
}

// PasswordHandler returns base64 value of the hashed password
func PasswordHandler(path string) http.Handler {
	log.Printf("Initiating Password hash handler %s", path)
	fn := func(w http.ResponseWriter, request *http.Request) {

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

		if request.Method != http.MethodPost {
			log.Printf("Requested method %s is not allowed", request.Method)
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)

		} else {
			b, err := ioutil.ReadAll(request.Body)
			if err != nil {
				log.Fatalln("Error when reading the request body", err)
			}

			// validate if the request body has required format
			payload := string(b)
			if strings.Index(payload, passwordPrefix) != 0 {
				http.Error(w, "Invalid request payload", http.StatusBadRequest)
				return
			}

			// split everything after prefix as one string
			payloadArray := strings.SplitAfterN(payload, passwordPrefix, 2)
			if len(payloadArray) <= 1 {
				http.Error(w, "Invalid request payload", http.StatusBadRequest)
				return
			}

			returnValue, err := passwordutil.GeneratePasswordHash(payloadArray[1]) // generate the hash
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			// OK, return the hashed password
			time.Sleep(5 * time.Second) // sleeping for 5 seconds
			setResponseCodeAndCustomMsg(w, http.StatusOK, returnValue)
		}
	}
	return http.HandlerFunc(fn)
}

// ShutdownHandler accepts the request to graefully shutdown the server; Returns OK 201 immediately
func ShutdownHandler(path string, done chan bool) http.Handler {
	log.Printf("Initiating ShutdownHandler for path %s", path)
	hanlder := func(w http.ResponseWriter, r *http.Request) {
		setResponseCodeAndDefaultMsg(w, http.StatusAccepted)
		done <- true
	}
	return http.HandlerFunc(hanlder)
}
