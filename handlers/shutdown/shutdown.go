package shutdownHandler

import (
	"context"
	"log"
	"net/http"
)

// GracefulShutdown waits for all the active requests to complete and shuts down the server gracefully
func GracefulShutdown(hs *http.Server, shutdownServer chan bool) {
	<-shutdownServer
	log.Println("Gracefully shutting down the server")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err := hs.Shutdown(ctx)
	if err != nil {
		log.Printf("Error when shutting down the server : %v", err)
	} else {
		log.Println("Graceful shutdown, server stopped")
	}
}
