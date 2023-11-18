package pg

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Masterminds/squirrel"
	"log/slog"
	"strings"

	"github.com/insan1a/exile/internal/models"
	"github.com/insan1a/exile/internal/storage"
	"github.com/insan1a/exile/internal/storage/person"
)

type Storage struct {
	db *sql.DB
}

// New creates a new user storage.
//
// If db is nil returns storage.ErrNilDB.
func New(db *sql.DB) (*Storage, error) {
	if db == nil {
		return nil, storage.ErrNilDB
	}

	return &Storage{db: db}, nil
}

// FindByID returns a person by given id
//
// If user not found returns person.ErrNotFound.
func (s *Storage) FindByID(ctx context.Context, id string) (*models.Person, error) {
	const query = `
	SELECT
		id,
		name,
		surname,
		patronymic,
		age,
		gender,
		nationality
	FROM person
	WHERE id = $1 AND is_deleted = FALSE
		`

	stmt, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("Storage.FindByID: %w", err)
	}
	defer stmt.Close()

	var p models.Person
	if err = stmt.QueryRowContext(ctx, id).
		Scan(&p.ID, &p.Name, &p.Surname, &p.Patronymic, &p.Age, &p.Gender, &p.Nationality); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("Storage.FindByID: %w", person.ErrNotFound)
		}

		return nil, fmt.Errorf("Storage.FindByID: %w", err)
	}

	return &p, nil
}

// Update updates a person.
//
// If person is nil returns person.ErrNilPerson.
// If person not found or the person ID is empty string returns person.ErrNotFound.
func (s *Storage) Update(ctx context.Context, p *models.Person) error {
	if p == nil {
		return fmt.Errorf("Storage.Update: %w", person.ErrNilPerson)
	}

	if p.ID == "" {
		return fmt.Errorf("Storage.Update: %w", person.ErrNotFound)
	}

	query := `UPDATE person SET `
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
		return nil
	}

	args = append(args, p.ID)
	query += fmt.Sprintf(
		"%s WHERE id = $%d RETURNING name, surname, patronymic, age, gender, nationality",
		strings.Join(queryParts, ", "),
		len(args),
	)

	stmt, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("Storage.Update: %w", err)
	}
	defer stmt.Close()

	if err = stmt.QueryRowContext(ctx, args...).
		Scan(&p.Name, &p.Surname, &p.Patronymic, &p.Age, &p.Gender, &p.Nationality); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("Storage.Update: %w", person.ErrNotFound)
		}

		return fmt.Errorf("Storage.Update: %w", err)
	}

	return nil
}

// Create creates a new person.
//
// The ID and CreatedOn must be filled by the database.
func (s *Storage) Create(ctx context.Context, p *models.Person) error {
	const query = `
	INSERT INTO person
		(name, surname, patronymic, age, gender, nationality)
	VALUES
		($1, $2, $3, $4, $5, $6)
	RETURNING id, name, surname, patronymic, age, gender, nationality
	`

	stmt, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("Storage.Create: %w", err)
	}
	defer stmt.Close()

	err = stmt.QueryRowContext(ctx, p.Name, p.Surname, p.Patronymic, p.Age, p.Gender, p.Nationality).
		Scan(&p.ID, &p.Name, &p.Surname, &p.Patronymic, &p.Age, &p.Gender, &p.Nationality)
	if err != nil {
		return fmt.Errorf("Storage.Create: %w", err)
	}

	return nil
}

// List returns a list of persons by given filter params.
//
// If param have unsupported type returns storage.ErrUnsupportedParamType.
// The param type can be int, float64, float32, uint, string.
func (s *Storage) List(ctx context.Context, filter *models.Filter) ([]models.Person, error) {
	query, args, err := filter.Query().
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("Storage.List: %w", err)
	}

	slog.Info("list query", "query", query, "args", args)

	stmt, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("Storage.List: %w", err)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, args...)
	if err != nil {
		return nil, fmt.Errorf("Storage.List: %w", err)
	}
	defer rows.Close()

	var people []models.Person
	for rows.Next() {
		var p models.Person
		if err = rows.Scan(&p.ID, &p.Name, &p.Surname, &p.Patronymic, &p.Age, &p.Gender, &p.Nationality); err != nil {
			return nil, fmt.Errorf("Storage.List: %w", err)
		}
		people = append(people, p)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("Storage.List: %w", err)
	}

	return people, nil
}

// Delete deletes a person by ID.
//
// Actually it sets is_delete field to true in database.
func (s *Storage) Delete(ctx context.Context, id string) error {
	const query = "UPDATE person SET is_deleted = TRUE WHERE id = $1"

	stmt, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("Storage.Delete: %w", err)
	}
	defer stmt.Close()

	if res, err := stmt.ExecContext(ctx, id); err != nil {
		return err
	} else {
		if count, err := res.RowsAffected(); err != nil || count == 0 {
			return fmt.Errorf("Storage.Delete: %w", person.ErrNotFound)
		}
	}

	return nil
}
