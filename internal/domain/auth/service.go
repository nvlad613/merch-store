package auth

import "context"

type Service interface {
	MakeAuth(user User, ctx context.Context) (JwtToken, error)
}
