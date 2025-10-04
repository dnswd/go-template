package server

import (
	"time"

	"github.com/go-chi/chi/middleware"
)

func (s *Server) setupMiddlewares() {
	s.router.Use(middleware.RequestID)
	s.router.Use(middleware.RealIP)
	s.router.Use(middleware.Logger)
	s.router.Use(middleware.Recoverer)
	s.router.Use(middleware.Timeout(60 * time.Second))
	// s.router.Use(RecoverMiddleware)
}

// func RecoverMiddleware(next http.Handler) http.Handler {
//     return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//         defer func() {
//             if err := recover(); err != nil {
//                 log.Printf("Panic recovered: %v\n%s", err, debug.Stack())

//                 RespondError(w,
//                     http.StatusInternalServerError,
//                     ErrCodeInternal,
//                     "An unexpected error occurred")
//             }
//         }()

//         next.ServeHTTP(w, r)
//     })
// }
