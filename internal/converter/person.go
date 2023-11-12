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
