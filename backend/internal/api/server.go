// internal/api/server.go
package api

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

// Server represents the HTTP server
type Server struct {
	router *gin.Engine
	server *http.Server
}

// NewServer creates a new server instance
func NewServer(router *gin.Engine, port string) *Server {
	return &Server{
		router: router,
		server: &http.Server{
			Addr:    ":" + port,
			Handler: router,
		},
	}
}

// Start starts the HTTP server with graceful shutdown
func (s *Server) Start() error {
	// Channel to listen for errors coming from the listener
	serverErrors := make(chan error, 1)

	// Start the server
	go func() {
		log.Printf("Server is running on port %s", s.server.Addr)
		serverErrors <- s.server.ListenAndServe()
	}()

	// Channel to listen for interrupt signals
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Block until we receive a signal or an error
	select {
	case err := <-serverErrors:
		return err
	case <-shutdown:
		log.Println("Shutting down server...")

		// Create a context with timeout for the graceful shutdown
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Attempt graceful shutdown
		err := s.server.Shutdown(ctx)
		if err != nil {
			// Force close if graceful shutdown fails
			log.Printf("Could not gracefully shutdown the server: %v\n", err)
			err = s.server.Close()
			if err != nil {
				return err
			}
		}
		log.Println("Server gracefully stopped")
	}

	return nil
}
