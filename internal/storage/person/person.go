package person

import (
	"context"
	"errors"

	"github.com/insan1a/exile/internal/models"
)

var (
	ErrNotFoundMany = errors.New("the people not found")
	ErrNotFound     = errors.New("the person not found")
	ErrNilPerson    = errors.New("the person could not be nil")
)

//go:generate go run github.com/vektra/mockery/v2@v2.20.2 --name Storage --output ./mocks --outpkg mocks
type Storage interface {
	FindByID(context.Context, string) (*models.Person, error)
	Update(context.Context, *models.Person) error
	Create(context.Context, *models.Person) error
	List(context.Context, *models.Filter) ([]models.Person, error)
	Delete(context.Context, string) error
}
