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
	gracefulShutdownPath = "/shutdown"
)

// main starts the server listening to different URL patterns and handles them
func main() {
	shutdownServer := make(chan bool, 1)

	mux := http.NewServeMux()
	server := &http.Server{
		Addr:    address,
		Handler: mux,
	}

	go func() {
		mux.Handle(rootPath, httphandler.RootHandler(rootPath))                                             // Application root handler
		mux.Handle(hashPasswordPath, httphandler.PasswordHandler(hashPasswordPath))                         // Password Hash handler
		mux.Handle(gracefulShutdownPath, httphandler.ShutdownHandler(gracefulShutdownPath, shutdownServer)) // Graceful shutdown handler
		server.ListenAndServe()                                                                             // make this a goroutine as its blocking?
	}()

	shutdown.GracefulShutdown(server, shutdownServer)
}
