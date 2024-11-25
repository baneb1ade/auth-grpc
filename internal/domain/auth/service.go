package auth

import (
	"auth-microserivce/internal/domain/models"
	"auth-microserivce/internal/lib/jwt"
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"time"
)

type Service struct {
	logger   *slog.Logger
	storage  Storage
	tokenTTL time.Duration
	secret   string
}

func NewService(logger *slog.Logger, storage Storage, tokenTTL time.Duration, secret string) *Service {
	return &Service{logger, storage, tokenTTL, secret}
}

func (s *Service) Login(ctx context.Context, username, password string) (string, error) {
	const op = "auth.service.Login"
	log := s.logger.With(slog.String("op", op))

	u, err := s.storage.GetOneByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", ErrInvalidUsernameOrPassword
		}
		log.Error(err.Error())
		return "", err
	}
	if err = bcrypt.CompareHashAndPassword(u.PassHash, []byte(password)); err != nil {
		return "", ErrInvalidUsernameOrPassword
	}

	token, err := jwt.NewToken(u, s.tokenTTL, s.secret)
	if err != nil {
		log.Error(err.Error())
	}
	return token, nil
}

func (s *Service) Register(ctx context.Context, email, username, password string) (string, error) {
	const op = "auth.service.Register"
	log := s.logger.With(slog.String("op", op))

	check, err := s.storage.IsEmailExists(ctx, email)
	if err != nil {
		return "", err
	}
	if check {
		return "", ErrEmailAlreadyExists
	}
	check, err = s.storage.IsUsernameExists(ctx, username)
	if err != nil {
		return "", err
	}
	if check {
		return "", ErrUsernameAlreadyExists
	}

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate password hash", err.Error())
	}
	id, err := s.storage.SaveOne(ctx, models.User{
		Email:    email,
		Username: username,
		PassHash: passHash,
	})
	if err != nil {
		log.Error("failed to save user", err.Error())
		return "", err
	}

	return id, nil
}
