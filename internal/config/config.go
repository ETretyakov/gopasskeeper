package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

// EnvMode is a string-like type to define application environment.
type EnvMode string

const (
	// LocalMode is a EnvMode type constant to define local environment.
	LocalMode EnvMode = "local"
	// DevMode is a EnvMode type constant to define dev environment.
	DevMode EnvMode = "dev"
	// ProdMode is a EnvMode type constant to define prod environment.
	ProdMode EnvMode = "prod"
)

// AppConfig is a configuration for app options.
type AppConfig struct {
	Env           EnvMode `env:"ENV"             envDefault:"dev"`
	LogLevel      string  `env:"LOG_LEVEL"       envDefault:"debug"`
	LogOutputFile string  `env:"LOG_OUTPUT_FILE" envDefault:"logs/stdout.log"`
	HealthHost    string  `env:"HEALTH_HOST"     envDefault:"127.0.0.1"`
	HealthPort    uint32  `env:"HEALTH_PORT"     envDefault:"8080"`
}

// LogLevel is a method to parse zerolog log level.
func (a *AppConfig) Level() zerolog.Level {
	if logLevel, err := zerolog.ParseLevel(a.LogLevel); err != nil {
		log.Warn().
			Str("logLevel", a.LogLevel).
			Msg("failed to parse log level: falling back to debug level")

		return zerolog.DebugLevel
	} else {
		return logLevel
	}
}

// LogOutputFile is a method to provide log file output.
func (a *AppConfig) OutputFile() string {
	return a.LogOutputFile
}

// HealthAddress is a string declaring http health address.
func (a *AppConfig) HealthAddress() string {
	return fmt.Sprintf("%s:%d", a.HealthHost, a.HealthPort)
}

// Config is configuration that aggregates all the others ones.
type Config struct {
	App        *AppConfig        `envPrefix:"APP_"`
	S3         *S3Config         `envPrefix:"S3_"`
	PostgreSQL *PostgreSQLConfig `envPrefix:"POSTGRES_"`
	Security   *SecurityConfig   `envPrefix:"SECURITY_"`
	Server     *ServerConfig     `envPrefix:"SERVER_"`
}

func (c *Config) Show() {
	msg := fmt.Sprintf(
		`APP: %+v S3: %+v PostgreSQL: %+v Security: %+v Server: %+v`,
		c.App, c.S3, c.PostgreSQL, c.Security, c.Server,
	)

	log.Debug().Msg(msg)
}

// New is a builder function to initiate Config structure.
func New() *Config {
	return &Config{
		App:        &AppConfig{},
		S3:         &S3Config{},
		PostgreSQL: &PostgreSQLConfig{},
		Security:   &SecurityConfig{},
		Server:     &ServerConfig{},
	}
}

// Option is a function type to declare signature for the config decorator.
type Option func(*Config)

// WithEnvFile is a decorator function to set up configuration file.
func WithEnvFile(envPath string, configType string) func(*Config) {
	return func(cfg *Config) {
		viper.SetConfigFile(envPath)
		viper.SetConfigType(configType)
		viper.AutomaticEnv()

		if err := viper.ReadInConfig(); err != nil {
			log.Error().Err(err).Msg("failed to load config from file")
		} else if err := viper.Unmarshal(cfg); err != nil {
			log.Error().Err(err).Msg("failed to parse config file")
		}
	}
}

func applyOptions(cfg *Config, options ...Option) {
	for i, option := range options {
		log.Info().Msg(fmt.Sprintf("applying config option %d/%d", i+1, len(options)))
		option(cfg)
	}
}

// LoadConfig is a builder function to parse Config structure from envs.
func LoadConfig(options ...Option) (*Config, error) {
	log.Info().Msg("reading configuration")
	cfg := New()

	if err := env.Parse(cfg); err != nil {
		return nil, errors.Wrap(err, "failed to load config")
	}

	applyOptions(cfg, options...)

	return cfg, nil
}
