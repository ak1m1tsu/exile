package models

import "fmt"

type Filter struct {
	Limit int `mapstructure:"limit" schema:"limit" validate:"omitempty,oneof=10 50 100"`
	Skip  int `mapstructure:"skip" schema:"skip" validate:"omitempty,gte=0"`

	Name        string `mapstructure:"name" schema:"name" validate:"omitempty,alpha"`
	Surname     string `mapstructure:"surname" schema:"surname" validate:"omitempty,alpha"`
	Patronymic  string `mapstructure:"patronymic" schema:"patronymic" validate:"omitempty,alpha"`
	Age         int    `mapstructure:"age" schema:"age" validate:"omitempty,gte=0,lte=150"`
	Gender      string `mapstructure:"gender" schema:"gender" validate:"omitempty,oneof=male female"`
	Nationality string `mapstructure:"nationality" schema:"nationality" validate:"omitempty,len=2"`
}

func (f Filter) String() string {
	return fmt.Sprintf(
		"limit: %d, skip: %d, name: %s, surname: %s, patronymic: %s, age: %d, gender: %s, nationality: %s",
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
