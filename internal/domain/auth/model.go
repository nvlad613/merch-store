package auth

type JwtToken = string

type User struct {
	Username     string
	Id           int
	PasswordHash []byte
	Coins        int
}

type Credentials struct {
	Username string
	Password string
}
