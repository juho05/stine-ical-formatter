package web

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/juho05/log"
)

type Server struct {
	httpServer http.Server
	renderer   *renderer
}

func NewServer(addr string) (*Server, error) {
	router := chi.NewMux()
	renderer, err := newRenderer()
	if err != nil {
		return nil, fmt.Errorf("new renderer: %w", err)
	}

	server := &Server{
		httpServer: http.Server{
			Addr:              addr,
			ReadTimeout:       30 * time.Second,
			ReadHeaderTimeout: 10 * time.Second,
			WriteTimeout:      30 * time.Second,
			ErrorLog:          log.NewStdLogger(log.WARNING),
			Handler:           router,
		},
		renderer: renderer,
	}

	registerMiddlewares(router)
	server.registerRoutes(router)

	return server, nil
}

func (s *Server) Listen() error {
	log.Infof("Listening on http://%s", s.httpServer.Addr)
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(timeout context.Context) error {
	log.Info("Shutting down...")
	return s.httpServer.Shutdown(timeout)
}
