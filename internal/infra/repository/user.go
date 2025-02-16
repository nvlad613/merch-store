package repository

import (
	"context"
	"database/sql"
	"errors"
	"merch-store/internal/domain"
	"merch-store/internal/domain/auth"
)

func (r *RepositoryImpl) GetUser(username string, ctx context.Context) (*auth.User, error) {
	var user User
	err := r.db.NewSelect().
		Model(&user).
		Where("name = ?", username).
		Scan(ctx)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.UserNotFoundError
		}

		return nil, err
	}
	userModel := user.ToModel()

	return &userModel, nil
}

func (r *RepositoryImpl) CreateUser(user auth.User, ctx context.Context) (int, error) {
	var id int
	_, err := r.db.NewInsert().
		Model(&user).
		Returning("id").
		Exec(ctx, &id)
	if err != nil {
		return 0, err
	}

	return id, nil
}
