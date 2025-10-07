package user

type User struct {
	ID           int
	Username     string
	PasswordHash string
	Name         string
	Role         string // "seller" or "buyer"
}