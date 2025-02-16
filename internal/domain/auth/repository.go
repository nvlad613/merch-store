package auth

import "context"

type Repository interface {
	HasUser(user User, ctx context.Context) (bool, error)
	CreateUser(user User, ctx context.Context) error
}
