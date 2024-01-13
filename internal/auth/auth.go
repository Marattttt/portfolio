package auth

import (
	"crypto/rand"
	"reflect"

	"github.com/Marattttt/portfolio/portfolio_back/internal/applog"
	"golang.org/x/crypto/bcrypt"
)

// Returns the hased password and the salt used to hash it
func HashSecret(secret []byte) ([]byte, []byte) {
	salt := getRandomSalt()
	hashed := hashSecret(secret, salt)

	return hashed, salt
}

func Validate(hashed, unhashed, salt []byte) bool {
	secret := append(unhashed, salt...)

	generated := hashSecret(secret, salt)

	return reflect.DeepEqual(secret, generated)
}

// Size of salt in bytes
const saltSize = 32

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
	salt := make([]byte, saltSize)

	if _, err := rand.Read(salt); err != nil {
		applog.Fatal(applog.Auth, err)
	}

	return salt
}
