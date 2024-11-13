package health

import "sync"

// Handler is a structure to declare liveness and readiness statuses.
type Handler struct {
	readinessMu    sync.RWMutex
	readinessState bool
	livenessMu     sync.RWMutex
	livenessState  bool
}

// SetupHandler is a builder function for Handler.
func SetupHandler() *Handler {
	return &Handler{}
}

// SetReadiness is a Handler method to set readiness status.
func (h *Handler) SetReadiness(state bool) {
	h.readinessMu.Lock()
	defer h.readinessMu.Unlock()
	h.readinessState = state
}

// SetLiveness is a Handler method to set liveness status.
func (h *Handler) SetLiveness(state bool) {
	h.livenessMu.Lock()
	defer h.livenessMu.Unlock()
	h.livenessState = state
}

// ReadinessState is a Handler method to retrieve readiness status.
func (h *Handler) ReadinessState() bool {
	h.readinessMu.RLock()
	defer h.readinessMu.RUnlock()
	return h.readinessState
}

// LivenessState is a Handler method to retrieve liveness status.
func (h *Handler) LivenessState() bool {
	h.livenessMu.RLock()
	defer h.livenessMu.RUnlock()
	return h.livenessState
}
