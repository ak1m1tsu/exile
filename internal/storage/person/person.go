package person

import (
	"context"
	"errors"

	"github.com/romankravchuk/effective-mobile-test-task/internal/models"
)

var (
	ErrNotFound  = errors.New("the person not found")
	ErrNilPerson = errors.New("the person could not be nil")
)

//go:generate go run github.com/vektra/mockery/v2@v2.33.3 --name Storage --output ./mocks --outpkg mocks
type Storage interface {
	FindByID(context.Context, string) (*models.Person, error)
	Update(context.Context, *models.Person) error
	Create(context.Context, models.Person) error
	List(context.Context, map[string]any) ([]models.Person, error)
}
