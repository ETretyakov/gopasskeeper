package grpc

import (
	"gopasskeeper/internal/closer"
	"gopasskeeper/internal/config"
	"gopasskeeper/internal/grpc/auth"
	"gopasskeeper/internal/grpc/interceptors"
	"gopasskeeper/internal/grpc/secretstore/accounts"
	"gopasskeeper/internal/grpc/secretstore/cards"
	"gopasskeeper/internal/grpc/secretstore/files"
	"gopasskeeper/internal/grpc/secretstore/notes"
	"gopasskeeper/internal/grpc/sync"
	"gopasskeeper/internal/lib/jwt"
	"gopasskeeper/internal/logger"

	accountsv1 "gopasskeeper/internal/grpc/secretstore/accounts/gen/accounts"
	cardsv1 "gopasskeeper/internal/grpc/secretstore/cards/gen/cards"
	filesv1 "gopasskeeper/internal/grpc/secretstore/files/gen/files"
	notesv1 "gopasskeeper/internal/grpc/secretstore/notes/gen/notes"
	syncv1 "gopasskeeper/internal/grpc/sync/gen/sync"

	"net"

	grpcLogging "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	grpcRecovery "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
)

// App is a grpc server application structure.
type App struct {
	gRPCServer *grpc.Server
	cfg        *config.ServerConfig
}

func accessibleRoles() map[string][]string {
	return map[string][]string{
		accountsv1.Accounts_Add_FullMethodName:       {"user"},
		accountsv1.Accounts_GetSecret_FullMethodName: {"user"},
		accountsv1.Accounts_Search_FullMethodName:    {"user"},
		accountsv1.Accounts_Remove_FullMethodName:    {"user"},
		cardsv1.Cards_Add_FullMethodName:             {"user"},
		cardsv1.Cards_GetSecret_FullMethodName:       {"user"},
		cardsv1.Cards_Search_FullMethodName:          {"user"},
		cardsv1.Cards_Remove_FullMethodName:          {"user"},
		notesv1.Notes_Add_FullMethodName:             {"user"},
		notesv1.Notes_GetSecret_FullMethodName:       {"user"},
		notesv1.Notes_Search_FullMethodName:          {"user"},
		notesv1.Notes_Remove_FullMethodName:          {"user"},
		filesv1.Files_Add_FullMethodName:             {"user"},
		filesv1.Files_GetSecret_FullMethodName:       {"user"},
		filesv1.Files_Search_FullMethodName:          {"user"},
		filesv1.Files_Remove_FullMethodName:          {"user"},
		syncv1.Sync_Get_FullMethodName:               {"user"},
	}
}

// New is a builder function to initiate App strucutre object.
func New(
	cfg *config.Config,
	authService auth.Auth,
	accountsService accounts.Accounts,
	cardsService cards.Cards,
	notesService notes.Notes,
	filesService files.Files,
	syncService sync.Sync,
) (*App, error) {
	loggingOpts := []grpcLogging.Option{}

	if cfg.App.Env != config.ProdMode {
		loggingOpts = []grpcLogging.Option{
			grpcLogging.WithLogOnEvents(
				grpcLogging.PayloadReceived,
				grpcLogging.PayloadSent,
			),
		}
	}

	recoveryOpts := []grpcRecovery.Option{
		grpcRecovery.WithRecoveryHandler(func(p interface{}) (err error) {
			logger.Error("recovered from panic", err)
			return status.Errorf(codes.Internal, "internal error")
		}),
	}

	jwtManager := jwt.NewJWTManager(cfg.Security.SignKey, cfg.Security.TokenTTL)

	interceptor := interceptors.NewAuthInterceptor(jwtManager, accessibleRoles())
	serverOptions := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			grpcRecovery.UnaryServerInterceptor(recoveryOpts...),
			grpcLogging.UnaryServerInterceptor(logger.InterceptorLogger(), loggingOpts...),
		),
		grpc.UnaryInterceptor(interceptor.Unary()),
		grpc.StreamInterceptor(interceptor.Stream()),
	}

	if cfg.Security.CertPath != "" {
		creds, err := credentials.NewServerTLSFromFile(cfg.Security.CertPath, cfg.Security.CertKeyPath)
		if err != nil {
			return nil, errors.Wrap(err, "cannot load TLS credentials")
		}

		serverOptions = append(serverOptions, grpc.Creds(creds))
	}

	gRPCServer := grpc.NewServer(serverOptions...)

	auth.Register(gRPCServer, authService)
	accounts.RegisterAccounts(gRPCServer, accountsService)
	cards.RegisterCards(gRPCServer, cardsService)
	notes.RegisterNotes(gRPCServer, notesService)
	files.RegisterFiles(gRPCServer, filesService)
	sync.RegisterSync(gRPCServer, syncService)

	return &App{
		gRPCServer: gRPCServer,
		cfg:        cfg.Server,
	}, nil
}

// Run is an App method to start running application server.
func (a *App) Run() error {
	const op = "grpcapp.Run"

	l, err := net.Listen("tcp", a.cfg.Address())
	if err != nil {
		return errors.Wrap(err, op)
	}

	logger.Info("grpc server started", "address", a.cfg.Address())

	if err := a.gRPCServer.Serve(l); err != nil {
		return errors.Wrap(err, op)
	}

	return nil
}

// MustRun is an App method to run application server with no returning error.
func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

// Start is an App method to start application server with a closer function.
func (a *App) Start() {
	go func() {
		a.MustRun()
	}()

	closer.Add(a.Close)
}

// Close is an App method to gracefully stop grpc server.
func (a *App) Close() error {
	logger.Info("stopping grpc server", "address", a.cfg.Address())
	a.gRPCServer.GracefulStop()
	return nil
}
