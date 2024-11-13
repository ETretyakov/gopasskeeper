package logger

import (
	"io"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	defaultFramesSkip int = 3
)

var (
	globalOutput io.Writer
	globalLogger zerolog.Logger
)

// LoggerConfig is an interface to describe logger config methods.
type LoggerConfig interface {
	Level() zerolog.Level
	OutputFile() string
}

// Init is an init function to setup global logger.
func Init(cfg LoggerConfig) {
	globalOutput = os.Stderr

	if outputFile := cfg.OutputFile(); outputFile != "" {
		runLogFile, err := os.OpenFile(
			outputFile,
			os.O_APPEND|os.O_CREATE|os.O_WRONLY,
			0664,
		)

		if err != nil {
			log.Error().Err(err).Msg("failed to open log output file")
		}

		multi := zerolog.MultiLevelWriter(globalOutput, runLogFile)
		globalLogger = zerolog.New(multi)
	} else {
		globalLogger = zerolog.New(globalOutput)
	}

	globalLogger = globalLogger.With().Stack().Logger()
	globalLogger = globalLogger.With().Timestamp().Logger()
	globalLogger = globalLogger.With().CallerWithSkipFrameCount(defaultFramesSkip).Logger()

	zerolog.SetGlobalLevel(zerolog.Level(cfg.Level()))
}

// Debug is a global function to write logs on debug level.
func Debug(msg string, kv ...any) {
	globalLogger.Debug().Fields(kv).Msg(msg)
}

// Info is a global function to write logs on info level.
func Info(msg string, kv ...any) {
	globalLogger.Info().Fields(kv).Msg(msg)
}

// Warn is a global function to write logs on warn level.
func Warn(msg string, kv ...any) {
	globalLogger.Warn().Fields(kv).Msg(msg)
}

// Error is a global function to write logs on error level.
func Error(msg string, err error, kv ...any) {
	globalLogger.Error().Fields(kv).Err(err).Msg(msg)
}
