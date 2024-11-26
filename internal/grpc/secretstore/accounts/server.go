package accounts

import (
	"context"
	"errors"

	accountsv1 "gopasskeeper/internal/grpc/secretstore/accounts/gen/accounts"

	"gopasskeeper/internal/domain/models"
	"gopasskeeper/internal/grpc/interceptors"

	"gopasskeeper/internal/services/secretstore/accounts"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type serverAPI struct {
	accountsv1.UnimplementedAccountsServer
	accounts Accounts
}

// Accounts is an interface for accounts storage.
type Accounts interface {
	Add(
		ctx context.Context,
		uid string,
		login string,
		server string,
		password string,
		meta string,
	) (*models.Message, error)
	GetSecret(
		ctx context.Context,
		uid string,
		accountID string,
	) (*models.AccountSecret, error)
	Search(
		ctx context.Context,
		uid string,
		schema *models.AccountSearchRequest,
	) (*models.AccountSearchResponse, error)
	Remove(
		ctx context.Context,
		uid string,
		accountID string,
	) (*models.Message, error)
}

// RegisterAccounts is a function to register Accounts server to grpc.
func RegisterAccounts(gRPCServer *grpc.Server, accounts Accounts) {
	accountsv1.RegisterAccountsServer(gRPCServer, &serverAPI{accounts: accounts})
}

// Add is a server api method to add accounts.
func (s *serverAPI) Add(
	ctx context.Context,
	in *accountsv1.AccountAddRequest,
) (*accountsv1.AccountAddResponse, error) {
	if (in.GetLogin() == "" || in.GetPassword() == "") && in.GetServer() == "" {
		return nil, status.Error(codes.InvalidArgument, "server is required and login or password is required")
	}

	uid, err := interceptors.ExtractUID(ctx)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "failed extract uid")
	}

	msg, err := s.accounts.Add(
		ctx,
		uid,
		in.GetLogin(),
		in.GetServer(),
		in.GetPassword(),
		in.GetMeta(),
	)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed add account")
	}

	return &accountsv1.AccountAddResponse{Status: msg.Status, Msg: msg.Msg}, nil
}

// GetSecret is a server api method to get account secret.
func (s *serverAPI) GetSecret(
	ctx context.Context,
	in *accountsv1.AccountSecretRequest,
) (*accountsv1.AccountSecretResponse, error) {
	if in.GetSecretId() == "" {
		return nil, status.Error(codes.InvalidArgument, "secret_id is required")
	}

	uid, err := interceptors.ExtractUID(ctx)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "failed extract uid")
	}

	accountSecret, err := s.accounts.GetSecret(ctx, uid, in.GetSecretId())
	if err != nil {
		if errors.Is(err, accounts.ErrAccountNotFound) {
			return nil, status.Error(codes.InvalidArgument, "failed to find account")
		}

		return nil, status.Error(codes.Internal, "failed to get secret")
	}

	return &accountsv1.AccountSecretResponse{
		Login:    accountSecret.Login,
		Server:   accountSecret.Server,
		Password: accountSecret.Password,
		Meta:     accountSecret.Meta,
	}, nil
}

// Search is a server api method to seatch account records.
func (s *serverAPI) Search(
	ctx context.Context,
	in *accountsv1.AccountSearchRequest,
) (*accountsv1.AccountSearchResponse, error) {
	if in.GetLimit() == 0 {
		return nil, status.Error(codes.InvalidArgument, "limit can't be 0")
	}

	uid, err := interceptors.ExtractUID(ctx)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "failed extract uid")
	}

	searchResponse, err := s.accounts.Search(
		ctx,
		uid,
		&models.AccountSearchRequest{
			Substring: in.GetSubstring(),
			Offset:    uint64(in.GetOffset()),
			Limit:     uint32(in.GetLimit()),
		},
	)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to search accounts")
	}

	items := []*accountsv1.AccountSearchItem{}
	for _, i := range searchResponse.Items {
		items = append(items, &accountsv1.AccountSearchItem{
			Id:     i.ID,
			Login:  i.Login,
			Server: i.Server,
		})
	}

	return &accountsv1.AccountSearchResponse{
		Count: int64(searchResponse.Count),
		Items: items,
	}, nil
}

// Remove is a server api method to remove an account record.
func (s *serverAPI) Remove(
	ctx context.Context,
	in *accountsv1.AccountRemoveRequest,
) (*accountsv1.AccountRemoveResponse, error) {
	if in.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "secret_id is required")
	}

	uid, err := interceptors.ExtractUID(ctx)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "failed extract uid")
	}

	msg, err := s.accounts.Remove(ctx, uid, in.GetId())
	if err != nil {
		if errors.Is(err, accounts.ErrAccountNotFound) {
			return nil, status.Error(codes.InvalidArgument, "failed to find account")
		}
		return nil, status.Error(codes.Internal, "failed to get account")
	}

	return &accountsv1.AccountRemoveResponse{
		Status: msg.Status,
		Msg:    msg.Msg,
	}, nil
}
