package psql

import (
	"auth-microserivce/internal/domain/auth"
	"auth-microserivce/internal/domain/models"
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"log/slog"
)

type Storage struct {
	Client auth.SQLClient
	Logger *slog.Logger
}

func NewStorage(client auth.SQLClient, logger *slog.Logger) *Storage {
	return &Storage{client, logger}
}

func (s *Storage) SaveOne(ctx context.Context, u models.User) (string, error) {
	const op = "db.psql.SaveOne"
	log := s.Logger.With(slog.String("op", op))

	q := `INSERT INTO "user" (username, email, password) VALUES ($1, $2, $3) RETURNING id`

	var id string
	if err := s.Client.QueryRow(ctx, q, u.Username, u.Email, u.PassHash).Scan(&id); err != nil {
		log.Error(op, "error", err)
		return "", err
	}
	return id, nil
}

func (s *Storage) GetOneByID(ctx context.Context, id string) (models.User, error) {
	const op = "db.psql.GetOneByID"
	log := s.Logger.With(slog.String("op", op))

	q := `SELECT id, username, email, password
          FROM "user"
          WHERE id = $1`

	var u models.User
	if err := s.Client.QueryRow(ctx, q, id).Scan(&u.ID, &u.Username, &u.Email, &u.PassHash); err != nil {
		log.Error(op, "error", err)
		return u, err
	}
	return u, nil
}

func (s *Storage) GetOneByUsername(ctx context.Context, username string) (models.User, error) {
	const op = "db.psql.GetOneByUsername"
	log := s.Logger.With(slog.String("op", op))

	q := `SELECT id, username, email, password
          FROM "user"
          WHERE username = $1`

	var u models.User
	if err := s.Client.QueryRow(ctx, q, username).Scan(&u.ID, &u.Username, &u.Email, &u.PassHash); err != nil {
		log.Error(op, "error", err)
		return u, err
	}
	return u, nil
}

func (s *Storage) GetOneByEmail(ctx context.Context, email string) (models.User, error) {
	const op = "db.psql.GetOneByEmail"
	log := s.Logger.With(slog.String("op", op))

	q := `SELECT id, username, email, password
          FROM "user"
          WHERE email = $1`

	var u models.User
	if err := s.Client.QueryRow(ctx, q, email).Scan(&u.ID, &u.Username, &u.Email, &u.PassHash); err != nil {
		log.Error(op, "error", err)

		return u, err
	}
	return u, nil
}

func (s *Storage) IsEmailExists(ctx context.Context, email string) (bool, error) {
	const op = "db.psql.IsEmailExists"
	log := s.Logger.With(slog.String("op", op))

	_, err := s.GetOneByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		log.Error(op, "error", err)
		return false, err
	}
	return true, nil
}

func (s *Storage) IsUsernameExists(ctx context.Context, username string) (bool, error) {
	const op = "db.psql.IsUsernameExists"
	log := s.Logger.With(slog.String("op", op))
	_, err := s.GetOneByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		log.Error(op, "error", err)
		return false, err
	}
	return true, nil
}
