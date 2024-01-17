package auth

import (
	"crypto/rand"
	"fmt"
	"unicode"

	"github.com/Marattttt/portfolio/portfolio_back/internal/applog"
	"golang.org/x/crypto/bcrypt"
)

const (
	// Max size for secret message
	SecretLength = 72 - SaltSize

	// Size of salt in bytes
	SaltSize = 32
)

type ErrDataTooLong struct{}

func (e ErrDataTooLong) Error() string {
	return fmt.Sprintf("Data too long")
}

// Returns the hased password and the salt used to hash it
func HashSecret(secret string) ([]byte, []byte) {
	salt := getRandomSalt()

	hashed := hashSecret([]byte(secret), salt)

	return hashed, salt
}

// Returns ErrDataTooLong if the data passed is too long
func Validate(hashed, salt []byte, unhashed string) error {
	if len(unhashed) > SecretLength || len(salt) > SaltSize {
		return ErrDataTooLong{}
	}

	secret := append([]byte(unhashed), salt...)

	if err := bcrypt.CompareHashAndPassword(hashed, secret); err != nil {
		return err
	}

	return nil
}

// Returns hashed password and the salt used
func hashSecret(secret, salt []byte) []byte {
	secret = append(secret, salt...)

	hashed, err := bcrypt.GenerateFromPassword(secret, bcrypt.DefaultCost)

	if err != nil {
		applog.Fatal(applog.Auth, err)
	}

	return hashed
}

func getRandomSalt() []byte {
	salt := make([]byte, SaltSize)

	if _, err := rand.Read(salt); err != nil {
		applog.Fatal(applog.Auth, err)
	}

	return salt
}

func getAsciiBytes(str string) *[]byte {
	bytes := make([]byte, len(str))

	for _, c := range str {
		if c > unicode.MaxASCII {
			return nil
		}

		bytes = append(bytes, byte(c))
	}

	return &bytes
}
