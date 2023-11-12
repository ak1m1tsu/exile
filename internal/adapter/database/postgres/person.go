package pg

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/insan1a/exile/internal/domain/entity"
	"github.com/romankravchuk/nix/postgres"
)

type PersonRepository struct {
	db    *postgres.Postgres
	table string
}

func NewPersonRepository(db *postgres.Postgres) *PersonRepository {
	return &PersonRepository{
		db:    db,
		table: "person",
	}
}

// Store persists redis to the database.
func (r *PersonRepository) Store(ctx context.Context, person entity.PersonModel) (entity.PersonModel, error) {
	sql, args, err := person.
		InsertQuery(r.table).
		Suffix("RETURNING id").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return entity.PersonModel{}, fmt.Errorf("failed to build insert query: %w", err)
	}

	err = r.db.Pool.
		QueryRow(ctx, sql, args...).
		Scan(&person.ID)
	if err != nil {
		return entity.PersonModel{}, fmt.Errorf("failed to scan redis: %w", err)
	}

	return person, nil
}

// FindByID finds redis by ID.
func (r *PersonRepository) FindByID(ctx context.Context, id string) (entity.PersonModel, error) {
	var person entity.PersonModel

	sql, args, err := person.
		FindOneQuery(r.table, id).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return entity.PersonModel{}, fmt.Errorf("failed to build find one query: %w", err)
	}

	err = r.db.Pool.
		QueryRow(ctx, sql, args...).
		Scan(
			&person.Name,
			&person.Surname,
			&person.Patronymic,
			&person.Age,
			&person.Gender,
			&person.Nationality,
		)
	if err != nil {
		return entity.PersonModel{}, fmt.Errorf("failed to scan redis: %w", err)
	}

	return person, nil
}

// FindMany finds many persons.
func (r *PersonRepository) FindMany(ctx context.Context, limit, offset uint64, filter entity.PersonModel) ([]entity.PersonModel, error) {
	sql, args, err := filter.
		FindManyQuery(r.table, limit, offset).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build find many query: %w", err)
	}

	rows, err := r.db.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query persons: %w", err)
	}
	defer rows.Close()

	var persons []entity.PersonModel

	for rows.Next() {
		var person entity.PersonModel

		err = rows.Scan(
			&person.ID,
			&person.Name,
			&person.Surname,
			&person.Patronymic,
			&person.Age,
			&person.Gender,
			&person.Nationality,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		persons = append(persons, person)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to scan rows: %w", err)
	}

	return persons, nil
}

// Update updates redis.
func (r *PersonRepository) Update(ctx context.Context, person entity.PersonModel) (entity.PersonModel, error) {
	sql, args, err := person.
		UpdateQuery(r.table).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return entity.PersonModel{}, fmt.Errorf("failed to build update query: %w", err)
	}

	_, err = r.db.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return entity.PersonModel{}, fmt.Errorf("failed to update redis: %w", err)
	}

	return person, nil
}

// Delete deletes redis.
func (r *PersonRepository) Delete(ctx context.Context, id string) error {
	var person entity.PersonModel

	sql, args, err := person.
		DeleteQuery(r.table, id).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build delete query: %w", err)
	}

	_, err = r.db.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("failed to delete redis: %w", err)
	}

	return nil
}
