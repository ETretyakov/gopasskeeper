package main

import (
	"context"
	"gopasskeeper/internal/clients/s3"
	"gopasskeeper/internal/closer"
	"gopasskeeper/internal/config"
	"gopasskeeper/internal/grpc"
	"gopasskeeper/internal/health"
	"gopasskeeper/internal/logger"
	"gopasskeeper/internal/repository"
	"gopasskeeper/internal/services/auth"
	"gopasskeeper/internal/services/secretstore/accounts"
	"gopasskeeper/internal/services/secretstore/cards"
	"gopasskeeper/internal/services/secretstore/files"
	"gopasskeeper/internal/services/secretstore/notes"
	"gopasskeeper/internal/services/sync"
	"gopasskeeper/internal/storage"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"
)

func main() {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer cancel()

	// load configuration
	cfg, err := config.LoadConfig(config.WithEnvFile("server-config.yaml", "yaml"))
	if err != nil {
		log.Error().Msg("failed to load config")
		return
	}

	cfg.Show()

	// init logger
	logger.Init(cfg.App)
	logger.Info("logger has been initiated")

	// health server
	healthServer := health.New(cfg.App)
	healthServer.Start()
	healthServer.SetLiveness(true)

	// init storage
	db, err := storage.NewPostgresDB(ctx, cfg.PostgreSQL)
	if err != nil {
		logger.Error("failed to connect to DB", err)
		return
	}

	repository := repository.New(db.DB)

	// init S3
	s3Client, err := s3.New(ctx, cfg.S3)
	if err != nil {
		logger.Error("failed to create s3 client", err)
		return
	}

	// setup grpc app
	authService := auth.New(
		cfg.Security,
		repository.Auth,
		repository.Auth,
	)

	accountsService, err := accounts.New(
		cfg.Security,
		repository.Accounts,
		repository.Sync,
	)
	if err != nil {
		logger.Error("failed to build account service", err)
		return
	}

	cardsService, err := cards.New(
		cfg.Security,
		repository.Cards,
		repository.Sync,
	)
	if err != nil {
		logger.Error("failed to build cards service", err)
		return
	}

	notesService, err := notes.New(
		cfg.Security,
		repository.Notes,
		repository.Sync,
	)
	if err != nil {
		logger.Error("failed to build notes service", err)
		return
	}

	filesService, err := files.New(
		cfg.Security,
		s3Client,
		repository.Files,
		repository.Sync,
	)
	if err != nil {
		logger.Error("failed to build files service", err)
		return
	}

	syncService := sync.New(
		repository.Sync,
	)

	app, err := grpc.New(
		cfg,
		authService,
		accountsService,
		cardsService,
		notesService,
		filesService,
		syncService,
	)
	if err != nil {
		logger.Error("failed to create app", err)
		return
	}

	app.Start()

	// graceful shutdown
	healthServer.SetReadiness(true)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	<-stop

	closer.CloseAll()
	logger.Info("gracefully stopped")
}
