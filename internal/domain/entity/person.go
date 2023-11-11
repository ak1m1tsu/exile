package entity

import (
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
	return PersonModel{
		ID:          p.ID,
		Name:        p.Name,
		Surname:     p.Surname,
		Patronymic:  p.Patronymic,
		Gender:      p.Gender,
		Nationality: p.Nationality,
		Age:         p.Age,
	}
}

func (p Person) Validate() error {
	return validator.New().Struct(p)
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

func (m PersonModel) ToEntity() Person {
	return Person{
		ID:          m.ID,
		Name:        m.Name,
		Surname:     m.Surname,
		Patronymic:  m.Patronymic,
		Gender:      m.Gender,
		Nationality: m.Nationality,
		Age:         m.Age,
	}
}

// InsertQuery returns insertBuilder with sql query:
//
//	INSERT INTO table (name, surname, patronymic, age, gender, nationality)
//	VALUES (?, ?, ?, ?, ?, ?)
func (p PersonModel) InsertQuery(table string) sq.InsertBuilder {
	builder := sq.StatementBuilder.
		Insert(table).
		Columns("name", "surname", "patronymic", "age", "gender", "nationality").
		Values(p.Name, p.Surname, p.Patronymic, p.Age, p.Gender, p.Nationality)

	return builder
}

// FindOneQuery returns findOneBuilder with sql query:
//
//	SELECT
//		name, surname, patronymic, age, gender, nationality
//	FROM table
//	WHERE id = ?
func (p PersonModel) FindOneQuery(table string, id string) sq.SelectBuilder {
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
func (p PersonModel) FindManyQuery(table string, limit, offset int) sq.SelectBuilder {
	builder := sq.StatementBuilder.
		Select("id", "name", "surname", "patronymic", "age", "gender", "nationality").
		From(table)

	if p.Name != "" {
		builder = builder.Where(sq.Like{"name": p.Name})
	}

	if p.Surname != "" {
		builder = builder.Where(sq.Like{"surname": p.Surname})
	}

	if p.Patronymic != "" {
		builder = builder.Where(sq.Like{"patronymic": p.Patronymic})
	}

	if p.Age != 0 {
		builder = builder.Where(sq.Eq{"age": p.Age})
	}

	if p.Gender != "" {
		builder = builder.Where(sq.Eq{"gender": p.Gender})
	}

	if p.Nationality != "" {
		builder = builder.Where(sq.Eq{"nationality": p.Nationality})
	}

	return builder.Limit(uint64(limit)).Offset(uint64(offset))
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
func (p PersonModel) UpdateQuery(table string) sq.UpdateBuilder {
	builder := sq.StatementBuilder.
		Update(table).
		Set("name", p.Name).
		Set("surname", p.Surname).
		Set("patronymic", p.Patronymic).
		Set("age", p.Age).
		Set("gender", p.Gender).
		Set("nationality", p.Nationality).
		Where(sq.Eq{"id": p.ID})

	return builder
}

// DeleteQuery returns deleteBuilder with sql query:
//
//	DELETE FROM table
//	WHERE id = ?
func (p PersonModel) DeleteQuery(table string, id string) sq.DeleteBuilder {
	builder := sq.StatementBuilder.
		Delete(table).
		Where(sq.Eq{"id": id})

	return builder
}
