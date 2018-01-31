package shutdownHandler

import (
	"context"
	"log"
	"net/http"
)

// GracefulShutdown waits for all the active requests to complete and shuts down the server gracefully
func GracefulShutdown(server *http.Server, shutdownServer chan bool) {

	<-shutdownServer // wait until shutdownServer channel receives the signal

	log.Println("Gracefully shutting down the server")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := server.Shutdown(ctx)
	if err != nil {
		log.Printf("Error when shutting down the server : %v", err)
	} else {
		log.Println("Server stopped")
	}
}
