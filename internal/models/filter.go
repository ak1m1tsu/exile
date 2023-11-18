package models

import (
	"fmt"
	sq "github.com/Masterminds/squirrel"
)

type Filter struct {
	Limit int `schema:"limit" validate:"omitempty,oneof=10 50 100"`
	Skip  int `schema:"skip" validate:"omitempty,gte=0"`

	Name        string `schema:"name" validate:"omitempty,alpha"`
	Surname     string `schema:"surname" validate:"omitempty,alpha"`
	Patronymic  string `schema:"patronymic" validate:"omitempty,alpha"`
	Age         int    `schema:"age" validate:"omitempty,gte=0,lte=150"`
	Gender      string `schema:"gender" validate:"omitempty,oneof=male female"`
	Nationality string `schema:"nationality" validate:"omitempty,len=2"`
}

func (f Filter) Query() sq.SelectBuilder {
	builder := sq.StatementBuilder.
		Select("id", "name", "surname", "patronymic", "age", "gender", "nationality").
		From("person")

	if f.Limit > 0 {
		builder = builder.Limit(uint64(f.Limit))
	} else {
		builder = builder.Limit(10)
	}

	if f.Skip > 0 {
		builder = builder.Offset(uint64(f.Skip))
	} else {
		builder = builder.Offset(0)
	}

	if f.Name != "" {
		builder = builder.Where(sq.Like{"name": f.Name})
	}

	if f.Surname != "" {
		builder = builder.Where(sq.Like{"surname": f.Surname})
	}

	if f.Patronymic != "" {
		builder = builder.Where(sq.Like{"patronymic": f.Patronymic})
	}

	if f.Age != 0 {
		builder = builder.Where(sq.Eq{"age": f.Age})
	}

	if f.Gender != "" {
		builder = builder.Where(sq.Eq{"gender": f.Gender})
	}

	if f.Nationality != "" {
		builder = builder.Where(sq.Eq{"nationality": f.Nationality})
	}

	return builder
}

func (f Filter) String() string {
	return fmt.Sprintf("filter-limit=%d-skip=%d-name=%s-surname=%s-patronymic=%s-age=%d-gender=%s-nationality=%s",
		f.Limit,
		f.Skip,
		f.Name,
		f.Surname,
		f.Patronymic,
		f.Age,
		f.Gender,
		f.Nationality,
	)
}
