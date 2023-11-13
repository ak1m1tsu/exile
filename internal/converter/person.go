package converter

import (
	"github.com/insan1a/exile/internal/controller/http/dto"
	"github.com/insan1a/exile/internal/domain/entity"
)

func PersonEntityToView(person entity.Person) dto.PersonView {
	return dto.PersonView(person)
}

func PersonEntitiesToViews(people []entity.Person) []dto.PersonView {
	views := make([]dto.PersonView, len(people))

	for i, person := range people {
		views[i] = PersonEntityToView(person)
	}

	return views
}

func FindPersonParamsToEntity(params dto.FindPersonParams) entity.Person {
	return entity.Person{
		Name:        params.Name,
		Surname:     params.Surname,
		Patronymic:  params.Patronymic,
		Gender:      params.Gender,
		Nationality: params.Nationality,
		Age:         params.Age,
	}
}

func CreatePersonDTOToEntity(person dto.CreatePersonDTO) entity.Person {
	return entity.Person{
		Name:       person.Name,
		Surname:    person.Surname,
		Patronymic: person.Patronymic,
	}
}

func UpdatePersonDTOToEntity(id string, person dto.UpdatePersonDTO) entity.Person {
	return entity.Person{
		ID:          id,
		Name:        person.Name,
		Surname:     person.Surname,
		Patronymic:  person.Patronymic,
		Gender:      person.Gender,
		Nationality: person.Nationality,
		Age:         person.Age,
	}
}
