package health

import (
	"context"
	"errors"
	"gopasskeeper/internal/closer"
	"gopasskeeper/internal/logger"
	"net/http"
)

// ServerConfig is an interface to declare health server configuration methods.
type ServerConfig interface {
	HealthAddress() string
}

// Server is a structure to describe http health server.
type Server struct {
	cfg    *ServerConfig
	srv    *http.Server
	health *Handler
}

// New is a builder function for health Server.
func New(cfg ServerConfig) *Server {
	srv := &http.Server{
		Addr: cfg.HealthAddress(),
	}

	healthHandler := SetupHandler()

	return &Server{
		cfg:    &cfg,
		srv:    srv,
		health: healthHandler,
	}
}

// Start is a Server method to start http health server.
func (s *Server) Start() {
	s.srv.Handler = s.setupRoutes()

	go func() {
		if err := s.srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			logger.Error("failed to start health server", err)
		}
	}()

	closer.Add(s.Close)
}

// Close is a Server function to gracefully stop health server.
func (s *Server) Close() error {
	if err := s.srv.Shutdown(context.Background()); err != nil {
		logger.Error("failed to close server", err)
	}

	logger.Info("health server is running")

	return nil
}

// SetLiveness is a Server method to set liveness status.
func (s *Server) SetLiveness(state bool) {
	s.health.SetLiveness(state)
}

// SetReadiness is a Server method to set readiness status.
func (s *Server) SetReadiness(state bool) {
	s.health.SetReadiness(state)
}
