package auth

import "context"

type Repository interface {
	GetUser(username string, ctx context.Context) (*User, error)
	CreateUser(user User, ctx context.Context) (int, error)
}
