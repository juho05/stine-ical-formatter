package web

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/didip/tollbooth/v7/limiter"
	"github.com/go-chi/chi/v5"
	"github.com/juho05/log"
)

type Server struct {
	httpServer http.Server
	renderer   *renderer
	limiter    *limiter.Limiter
	metrics    *Metrics
}

func NewServer(addr string) (*Server, error) {
	router := chi.NewMux()
	renderer, err := newRenderer()
	if err != nil {
		return nil, fmt.Errorf("new renderer: %w", err)
	}

	limit := limiter.New(nil).SetBurst(2).SetMax(2)
	if v, _ := strconv.ParseBool(os.Getenv("RATE_LIMIT_X_FORWARDED_FOR")); v {
		limit.SetIPLookups([]string{"X-Forwarded-For", "RemoteAddr"}).SetForwardedForIndexFromBehind(0)
	} else {
		limit.SetIPLookups([]string{"RemoteAddr"})
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
		limiter:  limit,
		metrics:  NewMetrics(),
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
