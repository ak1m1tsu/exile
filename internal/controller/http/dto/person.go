package dto

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/schema"
)

type CreatePersonDTO struct {
	Name       string `json:"name" validate:"required,alpha"`
	Surname    string `json:"surname" validate:"required,alpha"`
	Patronymic string `json:"patronymic" validate:"omitempty,alpha"`
}

// FromRequest fills CreatePersonDTO from request body and validates it.
func (dto *CreatePersonDTO) FromRequest(r *http.Request) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(dto); err != nil {
		return fmt.Errorf("failed to decode url query: %w", err)
	}

	if err := validator.New().Struct(dto); err != nil {
		return fmt.Errorf("failed to validate url query: %w", err)
	}

	return nil
}

type UpdatePersonDTO struct {
	Name        string `json:"name" validate:"required,alpha"`
	Surname     string `json:"surname" validate:"required,alpha"`
	Patronymic  string `json:"patronymic" validate:"omitempty,alpha"`
	Gender      string `json:"gender" validate:"required,oneof=male female"`
	Nationality string `json:"nationality" validate:"required,iso3166_1_alpha2"`
	Age         int    `json:"age" validate:"required,gt=0,lt=150"`
}

// FromRequest fills UpdatePersonDTO from request body and validates it.
func (p *UpdatePersonDTO) FromRequest(r *http.Request) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(p); err != nil {
		return fmt.Errorf("failed to decode url query: %w", err)
	}

	if err := validator.New().Struct(p); err != nil {
		return fmt.Errorf("failed to validate url query: %w", err)
	}

	return nil
}

type PersonView struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Surname     string `json:"surname"`
	Patronymic  string `json:"patronymic,omitempty"`
	Gender      string `json:"gender"`
	Nationality string `json:"nationality"`
	Age         int    `json:"age"`
}

type Pagination struct {
	Limit  uint64 `schema:"limit,required" validate:"required,oneof=10 30 50 100"`
	Offset uint64 `schema:"offset" validate:"omitempty,lte=100"`
}

type FindPersonParams struct {
	Name        string `schema:"name" validate:"omitempty,alpha"`
	Surname     string `schema:"surname" validate:"omitempty,alpha"`
	Patronymic  string `schema:"patronymic" validate:"omitempty,alpha"`
	Gender      string `schema:"gender" validate:"omitempty,oneof=male female"`
	Nationality string `schema:"nationality" validate:"omitempty,iso3166_1_alpha2"`
	Age         int    `schema:"age" validate:"omitempty,gt=0,lt=150"`
	Pagination
}

// FromRequest fills FindPersonParams from request url query and validates it.
func (fpp *FindPersonParams) FromRequest(r *http.Request) error {
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(false)

	if err := decoder.Decode(fpp, r.URL.Query()); err != nil {
		return fmt.Errorf("failed to decode url query: %w", err)
	}

	if err := validator.New().Struct(fpp); err != nil {
		return fmt.Errorf("failed to validate url query: %w", err)
	}

	return nil
}
