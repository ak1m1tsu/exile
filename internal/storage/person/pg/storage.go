package pg

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/romankravchuk/effective-mobile-test-task/internal/lib/errtools"
	"github.com/romankravchuk/effective-mobile-test-task/internal/models"
	"github.com/romankravchuk/effective-mobile-test-task/internal/storage"
	"github.com/romankravchuk/effective-mobile-test-task/internal/storage/person"
)

type Storage struct {
	db *sql.DB
}

// New creates a new user storage.
//
// If db is nil returns storage.ErrNilDB.
func New(db *sql.DB) (*Storage, error) {
	const op = "storage.user.pg.New"

	if db == nil {
		return nil, errtools.WithOperation(storage.ErrNilDB, op)
	}

	return &Storage{db: db}, nil
}

// FindByID returns a person by given id
//
// If user not found returns person.ErrNotFound.
func (s *Storage) FindByID(ctx context.Context, id string) (*models.Person, error) {
	const (
		op    = "storage.user.pg.Storage.FindByID"
		query = `
		SELECT
			id,
			name,
			surname,
			patronymic,
			age,
			gender,
			nationality,
			created_on
		FROM people
		WHERE id = $1
		`
	)

	stmt, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, errtools.WithOperation(err, op)
	}
	defer stmt.Close()

	var p *models.Person
	if err = stmt.QueryRowContext(ctx, id).Scan(p); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errtools.WithOperation(person.ErrNotFound, op)
		}

		return nil, errtools.WithOperation(err, op)
	}

	return p, nil
}

// Update updates a person.
//
// If person is nil returns person.ErrNilPerson.
// If person not found or the person ID is empty string returns person.ErrNotFound.
func (s *Storage) Update(ctx context.Context, p *models.Person) error {
	const op = "storage.user.pg.Storage.Update"

	if p == nil {
		return errtools.WithOperation(person.ErrNilPerson, op)
	}

	if p.ID == "" {
		return errtools.WithOperation(person.ErrNotFound, op)
	}

	query := `UPDATE people SET `
	queryParts := make([]string, 0, 6)
	args := make([]any, 0, 6)

	if p.Name != "" {
		queryParts = append(queryParts, fmt.Sprintf("name = $%d", len(args)+1))
		args = append(args, p.Name)
	}

	if p.Surname != "" {
		queryParts = append(queryParts, fmt.Sprintf("surname = $%d", len(args)+1))
		args = append(args, p.Surname)
	}

	if p.Patronymic != "" {
		queryParts = append(queryParts, fmt.Sprintf("patronymic = $%d", len(args)+1))
		args = append(args, p.Patronymic)
	}

	if p.Age != 0 {
		queryParts = append(queryParts, fmt.Sprintf("age = $%d", len(args)+1))
		args = append(args, p.Age)
	}

	if p.Gender != "" {
		queryParts = append(queryParts, fmt.Sprintf("gender = $%d", len(args)+1))
		args = append(args, p.Gender)
	}

	if p.Nationality != "" {
		queryParts = append(queryParts, fmt.Sprintf("nationality = $%d", len(args)+1))
		args = append(args, p.Nationality)
	}

	if len(args) == 0 {
		return errtools.WithOperation(person.ErrNotFound, op)
	}

	args = append(args, p.ID)
	query += fmt.Sprintf("%s WHERE id = $%d", strings.Join(queryParts, ", "), len(args)+1)

	stmt, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		return errtools.WithOperation(err, op)
	}
	defer stmt.Close()

	if err := stmt.QueryRowContext(ctx, args...).Scan(p); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errtools.WithOperation(person.ErrNotFound, op)
		}

		return errtools.WithOperation(err, op)
	}

	return nil
}

// Create creates a new person.
//
// The ID and CreatedOn must be filled by the database.
func (s *Storage) Create(ctx context.Context, p models.Person) error {
	const (
		op    = "storage.user.pg.Storage.Create"
		query = `
		INSERT INTO people
			(name, surname, patronymic, age, gender, nationality)
		VALUES
			($1, $2, $3, $4, $5, $6)
		`
	)

	stmt, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		return errtools.WithOperation(err, op)
	}
	defer stmt.Close()

	if _, err := stmt.ExecContext(ctx, p); err != nil {
		return errtools.WithOperation(err, op)
	}

	return nil
}

// List returns a list of persons by given filter params.
//
// If param have unsupported type returns storage.ErrUnsupportedParamType.
// The param type can be int, float64, float32, uint, string.
func (s *Storage) List(ctx context.Context, params map[string]any) ([]models.Person, error) {
	const op = "storage.user.pg.Storage.List"

	var values []any
	var where []string

	for k, v := range params {
		values = append(values, v)

		switch v.(type) {
		case int, float64, float32, uint:
			where = append(where, fmt.Sprintf("%s = $%d", k, len(values)))
		case string:
			where = append(where, fmt.Sprintf("%s ILIKE %s$%d%s", k, "'%", len(values), "%'"))
		default:
			return nil, errtools.WithOperation(storage.ErrUnsupportedParamType, op)
		}
	}

	query := fmt.Sprint(`
	SELECT
		id,
		name,
		surname,
		patronymic,
		age,
		gender,
		nationality,
		created_on
	FROM people
	WHERE `, strings.Join(where, " AND "), `
	ORDER BY created_on DESC`,
	)

	stmt, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, errtools.WithOperation(err, op)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, values...)
	if err != nil {
		return nil, errtools.WithOperation(err, op)
	}

	var people []models.Person
	for rows.Next() {
		var p models.Person
		if err = rows.Scan(&p); err != nil {
			break
		}
		people = append(people, p)
	}

	if closeErr := rows.Close(); closeErr != nil {
		return nil, errtools.WithOperation(closeErr, op)
	}

	if err != nil {
		return nil, errtools.WithOperation(err, op)
	}

	if err = rows.Err(); err != nil {
		return nil, errtools.WithOperation(err, op)
	}

	return people, nil
}
