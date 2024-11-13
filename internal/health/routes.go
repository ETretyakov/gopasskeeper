package health

import (
	"net/http"

	"github.com/rs/zerolog/log"
)

func (s *Server) setupRoutes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/liveness", func(w http.ResponseWriter, _ *http.Request) {
		if s.health.LivenessState() {
			w.WriteHeader(http.StatusOK)
			if _, err := w.Write([]byte("OK")); err != nil {
				log.Err(err).Msg("failed to respond to liveness request")
			}
			return
		}
	})
	mux.HandleFunc("/readiness", func(w http.ResponseWriter, _ *http.Request) {
		if s.health.ReadinessState() {
			w.WriteHeader(http.StatusOK)
			if _, err := w.Write([]byte("OK")); err != nil {
				log.Err(err).Msg("failed to respond to liveness request")
			}
			return
		}
	})

	return mux
}
