package logger

// GRPCLogger is a structure to describe grpc specific logger.
type GRPCLogger struct {
	ServiceName string
	Operator    string
}

// NewGRPCLogger is a builder function for GRPCLogger.
func NewGRPCLogger(serviceName string) *GRPCLogger {
	return &GRPCLogger{ServiceName: serviceName}
}

// WithOperator is a GRPCLogger method to add operator tag for logs.
func (l *GRPCLogger) WithOperator(op string) *GRPCLogger {
	return &GRPCLogger{
		ServiceName: l.ServiceName,
		Operator:    op,
	}
}

func (l *GRPCLogger) addTags(kv ...any) []any {
	if l.Operator != "" {
		kv = append(kv, "operator", l.Operator)
	}

	return append(kv, "service", l.ServiceName)
}

// Debug is a GRPCLogger method to write logs on debug level.
func (l *GRPCLogger) Debug(msg string, kv ...any) {
	globalLogger.Debug().Fields(l.addTags(kv...)).Msg(msg)
}

// Info is a GRPCLogger method to write logs on info level.
func (l *GRPCLogger) Info(msg string, kv ...any) {
	globalLogger.Info().Fields(l.addTags(kv...)).Msg(msg)
}

// Warn is a GRPCLogger method to write logs on warn level.
func (l *GRPCLogger) Warn(msg string, kv ...any) {
	globalLogger.Warn().Fields(l.addTags(kv...)).Msg(msg)
}

// Error is a GRPCLogger method to write logs on error level.
func (l *GRPCLogger) Error(msg string, err error, kv ...any) {
	globalLogger.Error().Fields(l.addTags(kv...)).Err(err).Msg(msg)
}
