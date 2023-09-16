package models

import "time"

type Person struct {
	ID          string    `db:"id" `
	Name        string    `db:"name"`
	Surname     string    `db:"surname"`
	Patronymic  string    `db:"patronymic"`
	Age         int       `db:"age"`
	Gender      string    `db:"gender"`
	Nationality string    `db:"nationality"`
	CreatedOn   time.Time `db:"created_on"`
}
