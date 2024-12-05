package auth

import (
	"context"
	"gopasskeeper/internal/config"
	"gopasskeeper/internal/domain/models"
	"gopasskeeper/internal/lib/jwt"
	"gopasskeeper/internal/logger"
	"gopasskeeper/internal/storage"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

var (
	// ErrInvalidCredentials is an error for invalid credentials cases.
	ErrInvalidCredentials = errors.New("invalid credentials")
)

// UserStorage is an interface to declare user storage methods.
type UserStorage interface {
	SaveUser(ctx context.Context, login string, passHash []byte) (uid string, err error)
	User(ctx context.Context, login string) (models.UserAuth, error)
}

// UserSaver is an interface to declare user saver methods.
type UserSaver interface {
	SaveUser(ctx context.Context, login string, passHash []byte) (uid string, err error)
}

// UserProvider is an interface to declare user provider methods.
type UserProvider interface {
	User(ctx context.Context, login string) (*models.UserAuth, error)
}

// Auth is a structure to describe auth service.
type Auth struct {
	log         *logger.GRPCLogger
	usrSaver    UserSaver
	usrProvider UserProvider
	secret      string
	tokenTTL    time.Duration
}

// New is a builder function for Auth.
func New(
	cfg *config.SecurityConfig,
	userSaver UserSaver,
	userProvider UserProvider,
) *Auth {
	return &Auth{
		log:         logger.NewGRPCLogger("auth"),
		usrSaver:    userSaver,
		usrProvider: userProvider,
		secret:      cfg.SignKey,
		tokenTTL:    cfg.TokenTTL,
	}
}

// RegisterNewUser is an Auth method to register new users within the auth service.
func (a *Auth) RegisterNewUser(
	ctx context.Context,
	login string,
	password string,
) (string, error) {
	const op = "Auth.RegisterNewUser"
	log := a.log.WithOperator(op)
	log.Info("registering user", "login", login)

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate password hash", err, "login", login)
		return "", errors.Wrap(err, op)
	}

	uuid, err := a.usrSaver.SaveUser(ctx, login, passwordHash)
	if err != nil {
		log.Error("failed to save user", err)
		return "", errors.Wrap(err, op)
	}

	return uuid, nil
}

// Login is an Auth method to login users within the auth service.
func (a *Auth) Login(
	ctx context.Context,
	login string,
	password string,
) (string, error) {
	const op = "Auth.Login"
	log := a.log.WithOperator(op)
	log.Info("attempting to login user", "login", login)

	user, err := a.usrProvider.User(ctx, login)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			log.Info("user not found", "login", login, "err", err)
			return "", errors.Wrap(ErrInvalidCredentials, op)
		}

		log.Error("failed to get user", err, "login", login)
		return "", errors.Wrap(err, op)
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		log.Error("invalid credentials", err, "login", login)
		return "", errors.Wrap(ErrInvalidCredentials, op)
	}

	jwtManage := jwt.NewJWTManager(a.secret, a.tokenTTL)
	token, err := jwtManage.Generate(user)
	if err != nil {
		log.Error("failed to generate token", err, "login", login)
		return "", errors.Wrap(err, op)
	}

	log.Info("user logged in successfully", "login", login)

	return token, nil
}
