package entity

import (
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/go-playground/validator/v10"
)

type Person struct {
	ID          string `validate:"omitempty,uuid"`
	Name        string `validate:"required,alpha"`
	Surname     string `validate:"required,alpha"`
	Patronymic  string `validate:"omitempty,alpha"`
	Gender      string `validate:"required,oneof=male female"`
	Nationality string `validate:"required,iso3166_2"`
	Age         int    `validate:"required,gt=0,lt=150"`
}

func (p Person) ToModel() PersonModel {
	return PersonModel(p)
}

func (p Person) Validate() error {
	if err := validator.New().Struct(p); err != nil {
		return fmt.Errorf("validation error: %w", err)
	}

	return nil
}

type PersonModel struct {
	ID          string
	Name        string
	Surname     string
	Patronymic  string
	Gender      string
	Nationality string
	Age         int
}

func (pm PersonModel) ToEntity() Person {
	return Person(pm)
}

// InsertQuery returns insertBuilder with sql query:
//
//	INSERT INTO table (name, surname, patronymic, age, gender, nationality)
//	VALUES (?, ?, ?, ?, ?, ?)
func (pm PersonModel) InsertQuery(table string) sq.InsertBuilder {
	builder := sq.StatementBuilder.
		Insert(table).
		Columns("name", "surname", "patronymic", "age", "gender", "nationality").
		Values(pm.Name, pm.Surname, pm.Patronymic, pm.Age, pm.Gender, pm.Nationality)

	return builder
}

// FindOneQuery returns findOneBuilder with sql query:
//
//	SELECT
//		name, surname, patronymic, age, gender, nationality
//	FROM table
//	WHERE id = ?
func (pm PersonModel) FindOneQuery(table string, id string) sq.SelectBuilder {
	builder := sq.StatementBuilder.
		Select("name", "surname", "patronymic", "age", "gender", "nationality").
		From(table).
		Where(sq.Eq{"id": id})

	return builder
}

// FindManyQuery returns findManyBuilder with sql query:
//
//	 SELECT
//			id, name, surname, patronymic, age, gender, nationality
//	 FROM
//			table
//	 WHERE
//			name LIKE ?
//			AND surname LIKE ?
//			AND patronymic LIKE ?
//			AND age = ?
//			AND gender = ?
//			AND nationality = ?
//	 LIMIT ? OFFSET ?
func (pm PersonModel) FindManyQuery(table string, limit, offset uint64) sq.SelectBuilder {
	builder := sq.StatementBuilder.
		Select("id", "name", "surname", "patronymic", "age", "gender", "nationality").
		From(table)

	if pm.Name != "" {
		builder = builder.Where(sq.Like{"name": pm.Name})
	}

	if pm.Surname != "" {
		builder = builder.Where(sq.Like{"surname": pm.Surname})
	}

	if pm.Patronymic != "" {
		builder = builder.Where(sq.Like{"patronymic": pm.Patronymic})
	}

	if pm.Age != 0 {
		builder = builder.Where(sq.Eq{"age": pm.Age})
	}

	if pm.Gender != "" {
		builder = builder.Where(sq.Eq{"gender": pm.Gender})
	}

	if pm.Nationality != "" {
		builder = builder.Where(sq.Eq{"nationality": pm.Nationality})
	}

	return builder.Limit(limit).Offset(offset)
}

// UpdateQuery returns updateBuilder with sql query:
//
//	UPDATE table
//	SET name = ?,
//	    surname = ?,
//	    patronymic = ?,
//	    age = ?,
//	    gender = ?,
//	    nationality = ?
//	WHERE id = ?
func (pm PersonModel) UpdateQuery(table string) sq.UpdateBuilder {
	builder := sq.StatementBuilder.
		Update(table).
		Set("name", pm.Name).
		Set("surname", pm.Surname).
		Set("patronymic", pm.Patronymic).
		Set("age", pm.Age).
		Set("gender", pm.Gender).
		Set("nationality", pm.Nationality).
		Where(sq.Eq{"id": pm.ID})

	return builder
}

// DeleteQuery returns deleteBuilder with sql query:
//
//	DELETE FROM table
//	WHERE id = ?
func (pm PersonModel) DeleteQuery(table string, id string) sq.DeleteBuilder {
	builder := sq.StatementBuilder.
		Delete(table).
		Where(sq.Eq{"id": id})

	return builder
}
