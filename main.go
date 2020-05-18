package main

import (
	"log"
	"net/http"

	"github.com/parthi-jaganathan/golang-assignment/handlers"
	"github.com/parthi-jaganathan/golang-assignment/handlers/shutdown"
)

const (
	address              = ":8080"
	rootPath             = "/"
	hashSequenceIDPath   = "/hash"
	hashPasswordByIDPath = "/hash/"
	gracefulShutdownPath = "/shutdown"
	statsPath            = "/stats"
)

// servePath invokes servermux handle
func servePath(pattern string, mux *http.ServeMux, handler http.Handler) {
	log.Printf("Server listening to path %s", pattern)
	mux.Handle(pattern, handler)
}

// main starts the server listening to different URL patterns and handles them
func main() {
	mux := http.NewServeMux()
	server := &http.Server{
		Addr:    address,
		Handler: mux,
	}

	shutdownServer := make(chan bool) // buffered channel

	servePath(rootPath, mux, httphandler.RootHandler())
	servePath(hashSequenceIDPath, mux, httphandler.GenerateRequestSequenceID())                         // Generate Request SequenceId
	servePath(hashPasswordByIDPath, mux, httphandler.GetPasswordHashBySequenceID(hashPasswordByIDPath)) // Get the Password Hash By request SequenceID
	servePath(gracefulShutdownPath, mux, httphandler.ShutdownGracefully(shutdownServer))                // Graceful shutdown handler
	servePath(statsPath, mux, httphandler.GetStats())                                                   // Stats path

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Println("something went wrong when starting server")
			shutdownServer <- true
		}
	}()

	// this will block until the channel receives shutdown server request
	shutdownHandler.GracefulShutdown(server, shutdownServer)
}
