package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"merch-store/config"
	"merch-store/internal/domain"
	"merch-store/pkg/jwtutil"
	"merch-store/pkg/timeprovider"
	"time"
)

type Service interface {
	MakeAuth(creds Credentials, ctx context.Context) (JwtToken, error)
}

type ServiceImpl struct {
	repository        Repository
	timerProvider     timeprovider.TimeProvider
	jwtExpirationTime time.Duration
	jwtSigningMethod  jwt.SigningMethod
	jwtSecretKey      []byte
}

type claims struct {
	Username string `json:"username"`
	UserID   int    `json:"user_id"`
	jwt.RegisteredClaims
}

func NewService(
	repository Repository,
	timeProvider timeprovider.TimeProvider,
	conf config.JwtAuth,
) (*ServiceImpl, error) {
	signingMethod, err := jwtutil.SigningMethodFromString(conf.Method)
	if err != nil {
		return nil, err
	}

	return &ServiceImpl{
		repository:        repository,
		timerProvider:     timeProvider,
		jwtSigningMethod:  signingMethod,
		jwtExpirationTime: time.Duration(conf.ExpSec) * time.Second,
		jwtSecretKey:      []byte(conf.Key),
	}, nil
}

func (s *ServiceImpl) MakeAuth(creds Credentials, ctx context.Context) (JwtToken, error) {
	user, err := s.repository.GetUser(creds.Username, ctx)
	if err != nil {
		if errors.Is(err, domain.UserNotFoundError) {
			userId, err := s.createUser(creds, ctx)
			if err != nil {
				return "", err
			}

			return s.genJwt(creds.Username, userId)
		}

		return "", err
	}

	if bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(creds.Password)) != nil {
		return "", WrongCredentialsError
	}

	return s.genJwt(user.Username, user.Id)
}

func (s *ServiceImpl) createUser(creds Credentials, ctx context.Context) (int, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(creds.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, fmt.Errorf("generate hash: %w", err)
	}

	id, err := s.repository.CreateUser(User{
		Username:     creds.Username,
		PasswordHash: hash,
	}, ctx)
	if err != nil {
		return 0, fmt.Errorf("create user: %w", err)
	}

	return id, nil
}

func (s *ServiceImpl) genJwt(username string, userId int) (JwtToken, error) {
	now := s.timerProvider.Now()
	expiration := now.Add(s.jwtExpirationTime)

	jwtClaims := claims{
		Username: username,
		UserID:   userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiration),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(s.jwtSigningMethod, jwtClaims)
	signedToken, err := token.SignedString(s.jwtSecretKey)
	if err != nil {
		return "", fmt.Errorf("signing jwt: %w", err)
	}

	return signedToken, nil
}
