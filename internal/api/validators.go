package api

import (
	"encoding/json"
	"io"

	"github.com/Marattttt/portfolio/portfolio_back/internal/auth"
)

type Validator interface {
	Validate() bool
}

// Returns a pointer to valid data, if the initial data was invalid, returns nil
func GetData[T Validator](r io.Reader) *T {
	var val T

	if err := json.NewDecoder(r).Decode(&val); err != nil {
		return nil
	}

	if !val.Validate() {
		return nil
	}
	return &val
}

func (r AuthRequest) Validate() bool {
	if r.Id <= 0 {
		return false
	}

	if len(r.Password) >= auth.SecretLength {
		return false
	}

	return true
}

func (r GuestRequest) Validate() bool {
	if len(r.Secret) >= auth.SecretLength {
		return false
	}

	return true
}
