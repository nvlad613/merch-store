package auth

type JwtToken = string

type User struct {
	Username     string
	Id           int
	PasswordHash []byte
}

type Credentials struct {
	Username string
	Password string
}
