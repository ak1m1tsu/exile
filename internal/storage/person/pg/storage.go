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
	const op = "storage.person.pg.New"

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
		op    = "storage.person.pg.Storage.FindByID"
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

	var p models.Person
	if err = stmt.QueryRowContext(ctx, id).
		Scan(&p.ID, &p.Name, &p.Surname, &p.Patronymic, &p.Age, &p.Gender, &p.Nationality, &p.CreatedOn); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errtools.WithOperation(person.ErrNotFound, op)
		}

		return nil, errtools.WithOperation(err, op)
	}

	return &p, nil
}

// Update updates a person.
//
// If person is nil returns person.ErrNilPerson.
// If person not found or the person ID is empty string returns person.ErrNotFound.
func (s *Storage) Update(ctx context.Context, p *models.Person) error {
	const op = "storage.person.pg.Storage.Update"

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
	query += fmt.Sprintf(
		"%s WHERE id = $%d RETURNING name, surname, patronymic, age, gender, nationality, created_on",
		strings.Join(queryParts, ", "),
		len(args),
	)

	stmt, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		return errtools.WithOperation(err, op)
	}
	defer stmt.Close()

	if err := stmt.QueryRowContext(ctx, args...).
		Scan(&p.Name, &p.Surname,
			&p.Patronymic, &p.Age, &p.Gender,
			&p.Nationality, &p.CreatedOn,
		); err != nil {
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
		op    = "storage.person.pg.Storage.Create"
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

	if _, err := stmt.ExecContext(ctx, p.Name, p.Surname, p.Patronymic, p.Age, p.Gender, p.Nationality); err != nil {
		return errtools.WithOperation(err, op)
	}

	return nil
}

// List returns a list of persons by given filter params.
//
// If param have unsupported type returns storage.ErrUnsupportedParamType.
// The param type can be int, float64, float32, uint, string.
func (s *Storage) List(ctx context.Context, filter *models.Filter) ([]models.Person, error) {
	const op = "storage.person.pg.Storage.List"

	var (
		where              []string
		args               int
		query, limit, skip string
	)

	if filter == nil {
		filter = &models.Filter{Limit: 10, Skip: 0}
	}
	if filter.Name != "" {
		args += 1
		where = append(where, fmt.Sprintf("name ILIKE '%s%s%s'", "%", filter.Name, "%"))
	}
	if filter.Surname != "" {
		args += 1
		where = append(where, fmt.Sprintf("surname ILIKE '%s%s%s'", "%", filter.Surname, "%"))
	}
	if filter.Patronymic != "" {
		args += 1
		where = append(where, fmt.Sprintf("patronymic ILIKE '%s%s%s'", "%", filter.Patronymic, "%"))
	}
	if filter.Age != 0 {
		args += 1
		where = append(where, fmt.Sprintf("age = %d", filter.Age))
	}
	if filter.Gender != "" {
		args += 1
		where = append(where, fmt.Sprintf("gender = '%s'", filter.Gender))
	}
	if filter.Nationality != "" {
		args += 1
		where = append(where, fmt.Sprintf("nationality = '%s'", filter.Nationality))
	}
	if filter.Limit == 0 {
		filter.Limit = 10
	}

	limit = fmt.Sprintf("LIMIT %d", filter.Limit)
	skip = fmt.Sprintf("OFFSET %d", filter.Skip)

	if args > 0 {
		query = fmt.Sprint(`SELECT id, name, surname, patronymic, age, gender, nationality, created_on
		FROM people
		WHERE `, strings.Join(where, " AND "), `
		ORDER BY created_on DESC `, limit, ` `, skip)
	} else {
		query = fmt.Sprint(`SELECT id, name, surname, patronymic, age, gender, nationality, created_on
		FROM people
		ORDER BY created_on DESC `, limit, ` `, skip)
	}

	stmt, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, errtools.WithOperation(err, op)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		return nil, errtools.WithOperation(err, op)
	}

	var people []models.Person
	for rows.Next() {
		var p models.Person
		if err = rows.Scan(
			&p.ID, &p.Name, &p.Surname,
			&p.Patronymic, &p.Age, &p.Gender,
			&p.Nationality, &p.CreatedOn,
		); err != nil {
			break
		}
		people = append(people, p)
	}

	if len(people) == 0 {
		return nil, errtools.WithOperation(person.ErrNotFoundMany, op)
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

// Delete deletes a person by ID.
//
// Actually it sets is_delete field to true in database.
func (s *Storage) Delete(ctx context.Context, id string) error {
	const (
		op    = "storage.person.pg.Storage.Delete"
		query = "UPDATE people SET is_deleted = true WHERE id = $1"
	)

	stmt, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		return errtools.WithOperation(err, op)
	}
	defer stmt.Close()

	if res, err := stmt.ExecContext(ctx, id); err != nil {
		return errtools.WithOperation(err, op)
	} else {
		if count, err := res.RowsAffected(); err != nil || count == 0 {
			return errtools.WithOperation(person.ErrNotFound, op)
		}
	}

	return nil
}
