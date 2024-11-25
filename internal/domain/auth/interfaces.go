package auth

import (
	"auth-microserivce/internal/domain/models"
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type SQLClient interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Begin(ctx context.Context) (pgx.Tx, error)
}

type Storage interface {
	SaveOne(ctx context.Context, user models.User) (string, error)
	GetOneByID(ctx context.Context, id string) (models.User, error)
	GetOneByUsername(ctx context.Context, username string) (models.User, error)
	GetOneByEmail(ctx context.Context, email string) (models.User, error)
	IsEmailExists(ctx context.Context, email string) (bool, error)
	IsUsernameExists(ctx context.Context, username string) (bool, error)
}
