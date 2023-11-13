package dto

type PersonEvent struct {
	Name       string
	Surname    string
	Patronymic string
}

type PersonErrorEvent struct {
	PersonEvent
	Error string
}
