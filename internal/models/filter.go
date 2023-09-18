package models

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
