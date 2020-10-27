package account

import (
	"context"
	"database/sql"
	"errors"
	"github.com/go-kit/kit/log"
)

var RepoErr = errors.New("Unable to handle repo request")

type repo struct {
	db     *sql.DB
	logger log.Logger
}

func (repo *repo) CreateUser(ctx context.Context, user User) error {
	query := "INSERT INTO users (id, email, password) VALUES (?, ?, ?)"
	if user.Email == "" || user.Password == "" {
		return RepoErr
	}

	_ , err := repo.db.ExecContext(ctx, query, user.ID, user.Email, user.Password)
	if err != nil {
		return err
	}
	return nil
}

func (repo *repo) GetUser(ctx context.Context, id string) (string, error) {
	var email string
	err := repo.db.QueryRow("SELECT email FROM users WHERE id = ?", id).Scan(&email)
	if err != nil {
		repo.logger.Log("Error from repo",id)
		return "", RepoErr
	}
	repo.logger.Log("Email received", email)
	return email, nil
}

func NewRepo(db *sql.DB, logger log.Logger) Repository {
	return &repo{
		db:     db,
		logger: log.With(logger, "repo", "sql"),
	}
}
