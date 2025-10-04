// server/server.go
package server

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/dnswd/arus/health"
	"github.com/dnswd/arus/user"
	"github.com/go-chi/chi"
)

type Server struct {
	router      chi.Router
	httpServer  *http.Server
	healthCheck *health.HttpHandler
	userHandler *user.HTTPHandler
	// orderHandler *order.Handler
}

func New(
	healthCheck *health.HttpHandler,
	userHandler *user.HTTPHandler,
	// orderHandler *order.Handler,
) *Server {
	r := chi.NewRouter()

	s := &Server{
		router:      r,
		healthCheck: healthCheck,
		userHandler: userHandler,
		// orderHandler: orderHandler,
	}

	s.setupMiddlewares()
	s.setupRoutes()

	return s
}

func (s *Server) Start(ctx context.Context, addr string) error {
	s.httpServer = &http.Server{
		Addr:         addr,
		Handler:      s.router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("Server starting on %s", addr)

	// Start in goroutine so it doesn't block
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	log.Println("Server shutting down...")
	return s.httpServer.Shutdown(ctx)
}
