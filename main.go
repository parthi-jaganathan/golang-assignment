package main

import (
	"net/http"

	"./handlers"
	"./handlers/shutdown"
)

const (
	address              = ":8080"
	rootPath             = "/"
	hashPasswordPath     = "/hash"
	hashPasswordPathID   = "/hash/"
	gracefulShutdownPath = "/shutdown"
	statsPath            = "/stats"
)

// main starts the server listening to different URL patterns and handles them
func main() {
	mux := http.NewServeMux()
	server := &http.Server{
		Addr:    address,
		Handler: mux,
	}

	shutdownServer := make(chan bool) // buffered channel
	go shutdownHandler.GracefulShutdown(server, shutdownServer)

	mux.Handle(rootPath, httphandler.RootHandler(rootPath))                                             // Application root handler
	mux.Handle(hashPasswordPath, httphandler.PasswordHandler(hashPasswordPath))                         // Password Hash handler
	mux.Handle(hashPasswordPathID, httphandler.PasswordHandlerID(hashPasswordPathID))                   // Password Hash handler
	mux.Handle(gracefulShutdownPath, httphandler.ShutdownHandler(gracefulShutdownPath, shutdownServer)) // Graceful shutdown handler
	mux.Handle(statsPath, httphandler.GetStats(statsPath))                                              // Stats path
	server.ListenAndServe()                                                                             // make this a goroutine as its blocking?
}
