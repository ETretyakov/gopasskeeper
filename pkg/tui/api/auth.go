package api

import (
	"context"
	"gopasskeeper/internal/clients/secretstore"
	ssov1 "gopasskeeper/internal/grpc/auth/gen/sso"
	"gopasskeeper/internal/logger"
	"time"

	"github.com/pkg/errors"
)

func (a *API) refreshToken() error {
	if a.client == nil {
		return nil
	}

	credentials := a.states.GetCredentials()
	loginResp, err := a.client.AuthAPI.Login(context.Background(), &ssov1.LoginRequest{
		Login:    credentials.Login,
		Password: credentials.Password,
	})
	if err != nil {
		return errors.Wrap(err, "failed to refresh token")
	}

	a.states.SetToken(loginResp.GetToken())

	return nil
}

func (a *API) Login(endpoint, login, password string) error {
	ctx := context.Background()
	log := logger.NewGRPCLogger("tui")

	client, err := secretstore.New(
		a.cfg,
		ctx,
		log,
		endpoint,
		time.Minute,
		3,
	)
	if err != nil {
		log.Error("failed to create client", err)
		return errors.Wrap(err, "failed to create client")
	}

	a.SetClient(client)
	a.states.SetCredentials(login, password)

	loginResp, err := client.AuthAPI.Login(ctx, &ssov1.LoginRequest{
		Login:    login,
		Password: password,
	})
	if err != nil {
		_, err := client.AuthAPI.Register(ctx, &ssov1.RegisterRequest{
			Login:    login,
			Password: password,
		})
		if err != nil {
			return errors.Wrap(err, "failed to register")
		}

		loginResp, err = client.AuthAPI.Login(ctx, &ssov1.LoginRequest{
			Login:    login,
			Password: password,
		})
		if err != nil {
			return errors.Wrap(err, "failed to login after register")
		}
	}

	a.states.SetToken(loginResp.GetToken())

	go func() {
		ticker := time.NewTicker(time.Minute)
		for {
			select {
			case <-ticker.C:
				if err := a.refreshToken(); err != nil {
					log.Error("failed to refresh token", err)
					return
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	return nil
}
