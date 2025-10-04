package server

import (
	"github.com/go-chi/chi"
)

func (s *Server) setupRoutes() {
	// Health checks (no auth needed)
	s.router.Get("/health/live", s.healthCheck.Liveness)
	s.router.Get("/health/ready", s.healthCheck.Readiness)

	// API routes
	s.router.Route("/api/v1", func(r chi.Router) {
		// TODO protected routes
		// User routes
		r.Route("/users", func(r chi.Router) {
			r.Post("/", s.userHandler.Create)
			r.Get("/{id}", s.userHandler.Get)
			r.Delete("/{id}", s.userHandler.Delete)
		})

		// // Order routes
		// r.Route("/orders", func(r chi.Router) {
		// 	r.Post("/", s.orderHandler.Create)
		// 	r.Get("/{id}", s.orderHandler.Get)
		// })
	})
}
