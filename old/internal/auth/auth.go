package auth

import (
	"fmt"
)

const (
	// Max size for secret message
	SecretLength = 72 - SaltSize

	// Size of salt in bytes
	SaltSize = 32
)

var (
	ErrUserExists    = fmt.Errorf("User already exists")
	ErrUserNotExists = fmt.Errorf("User does not exist")
)

type LoginData struct {
	Name     string
	Password string
}

type Manager interface {
	Register(LoginData) error
	Authenticate(LoginData) error
}

type ConnectableManager interface {
	Manager
	Connect() error
	Close() error
}
