package secretstore

import (
	"context"
	"time"

	ssov1 "gopasskeeper/internal/grpc/auth/gen/sso"
	accountsv1 "gopasskeeper/internal/grpc/secretstore/accounts/gen/accounts"
	cardsv1 "gopasskeeper/internal/grpc/secretstore/cards/gen/cards"
	filesv1 "gopasskeeper/internal/grpc/secretstore/files/gen/files"
	notesv1 "gopasskeeper/internal/grpc/secretstore/notes/gen/notes"
	syncv1 "gopasskeeper/internal/grpc/sync/gen/sync"
	"gopasskeeper/internal/logger"

	grpclog "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	grpcretry "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
)

type ClientConfig interface {
	CertPath() string
}

// Client is a structure to declare grpc client.
type Client struct {
	AuthAPI     ssov1.AuthClient
	AccountsAPI accountsv1.AccountsClient
	CardsAPI    cardsv1.CardsClient
	NotesAPI    notesv1.NotesClient
	FilesAPI    filesv1.FilesClient
	SyncAPI     syncv1.SyncClient
	logger      *logger.GRPCLogger
}

// New is a builder method for Client.
func New(
	cfg ClientConfig,
	ctx context.Context,
	log *logger.GRPCLogger,
	addr string,
	timeout time.Duration,
	retriesCount int,
) (*Client, error) {
	const op = "grpc.New"

	transportCreds := grpc.WithTransportCredentials(insecure.NewCredentials())
	if cfg.CertPath() != "" {
		creds, err := applyTLS(cfg.CertPath())
		if err != nil {
			return nil, errors.Wrap(err, "failed to apply tls")
		}

		transportCreds = grpc.WithTransportCredentials(creds)
	}

	retryOpts := []grpcretry.CallOption{
		grpcretry.WithCodes(codes.NotFound, codes.Aborted, codes.DeadlineExceeded),
		grpcretry.WithMax(uint(retriesCount)),
		grpcretry.WithPerRetryTimeout(timeout),
	}

	logOpts := []grpclog.Option{
		grpclog.WithLogOnEvents(grpclog.PayloadReceived, grpclog.PayloadSent),
	}

	cc, err := grpc.NewClient(
		addr,
		transportCreds,
		grpc.WithChainUnaryInterceptor(
			grpclog.UnaryClientInterceptor(logger.InterceptorLogger(), logOpts...),
			grpcretry.UnaryClientInterceptor(retryOpts...),
		),
	)
	if err != nil {
		log.Error("failed to build auth client", err, "service", op)
		return nil, errors.Wrap(err, "failed to build auth client")
	}

	authAPI := ssov1.NewAuthClient(cc)
	accountsAPI := accountsv1.NewAccountsClient(cc)
	cardsAPI := cardsv1.NewCardsClient(cc)
	notesAPI := notesv1.NewNotesClient(cc)
	filesAPI := filesv1.NewFilesClient(cc)
	syncAPI := syncv1.NewSyncClient(cc)

	return &Client{
		AuthAPI:     authAPI,
		AccountsAPI: accountsAPI,
		CardsAPI:    cardsAPI,
		NotesAPI:    notesAPI,
		FilesAPI:    filesAPI,
		SyncAPI:     syncAPI,
		logger:      log,
	}, nil
}
