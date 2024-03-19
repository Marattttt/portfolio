package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/Marattttt/portfolio/portfolio_back/internal/models"
)

type UsersRepository interface {
	// Returns nil if the entity is not found or was deleted
	Get(context.Context, int) (*models.User, error)

	// The newU parameter will have its id updated if the operation is successful
	//
	// If the id of the entity is non-zero, a check for existing entity is performed,
	// and if one exists ErrEntityExists is returned
	Create(context.Context, *models.User) error

	// The newU parameter will have its id updated if the operation is successful
	//
	// If entity does not exist, a EntityNotExistsError is returned
	Update(context.Context, *models.User) error

	// Returns optional nil if error or valid int ptr of the deleted id
	Delete(context.Context, int) error
}

// Should be used only if the existence is an error, not to indicate that something simply exists
//
// Contains id of the entity
type EntityExistsError struct {
	ID int
}

// Error interface for ErrEntityExists
func (e EntityExistsError) Error() string {
	return fmt.Sprintf("Entity id %d already exists", e.ID)
}

// Contains id of the entity
type EntityNotExistsError struct {
	ID int
}

// Error interface for ErrEntityNotExists
func (e EntityNotExistsError) Error() string {
	return fmt.Sprintf("Entity id %d does not exist", e.ID)
}

var (
	ErrNoEffect = errors.New("No rows affected")
)
