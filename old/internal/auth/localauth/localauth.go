package localauth

import (
	"crypto/rand"
	"fmt"

	"github.com/Marattttt/portfolio/portfolio_back/internal/auth"
	"github.com/Marattttt/portfolio/portfolio_back/internal/users"
	"golang.org/x/crypto/bcrypt"
)

type Manager struct {
	repo *users.Users
}

func NewManager(repository *users.Users) Manager {
	return Manager{
		repo: repository,
	}
}

func (m Manager) Register(login auth.LoginData) error {
	user, err := m.repo.GetName(login.Name)

	salt, err := getRandomSalt()
	if err != nil {
		return err
	}

	hashed, err := hashSecret([]byte(login.Password), salt)
	if err != nil {
		return err
	}
}

var DataTooLongError = fmt.Errorf("Provided data was too long")

// Returns the hased password and the salt used to hash it
func HashSecret(secret string) ([]byte, []byte, error) {
	salt, err := getRandomSalt()
	if err != nil {
		return nil, nil, err
	}

	hashed, err := hashSecret([]byte(secret), salt)
	if err != nil {
		return nil, nil, err
	}

	return hashed, salt, nil
}

// If data is too long returns DataTooLongError
func Validate(hashed, salt []byte, unhashed string) error {
	if len(unhashed) > SecretLength || len(salt) > SaltSize {
		return DataTooLongError{}
	}

	secret := append([]byte(unhashed), salt...)

	if err := bcrypt.CompareHashAndPassword(hashed, secret); err != nil {
		return fmt.Errorf("while comparing password: %w", err)
	}

	return nil
}

// Returns hashed password and the salt used
func hashSecret(secret, salt []byte) ([]byte, error) {
	secret = append(secret, salt...)

	hashed, err := bcrypt.GenerateFromPassword(secret, bcrypt.DefaultCost)

	if err != nil {
		return nil, fmt.Errorf("generating new password: %w", err)
	}

	return hashed, nil
}

func getRandomSalt() ([]byte, error) {
	salt := make([]byte, SaltSize)

	if _, err := rand.Read(salt); err != nil {
		return nil, fmt.Errorf("while generating new random salt: %w", err)
	}

	return salt, nil
}
