package v1

import (
	"context"
	"errors"
	"github.com/insan1a/exile/internal/converter"
	"net/http"

	chi "github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/schema"
	"github.com/insan1a/exile/internal/controller/http/dto"
	"github.com/insan1a/exile/internal/controller/http/response"
	"github.com/insan1a/exile/internal/domain/entity"
	log "github.com/romankravchuk/nix/log/zerolog"
)

const (
	personEndpoint   = "/person"
	personIDEndpoint = "/{person_id}"
	personIDParam    = "person_id"
)

type PersonService interface {
	Store(ctx context.Context, person entity.Person) (entity.Person, error)
	FindByID(ctx context.Context, id string) (entity.Person, error)
	FindMany(ctx context.Context, page, limit uint64, filter entity.Person) ([]entity.Person, error)
	Update(ctx context.Context, person entity.Person) (entity.Person, error)
	Delete(ctx context.Context, id string) error
}

type PersonHandler struct {
	service PersonService
	logger  *log.Logger
}

func MountPersonHandler(router chi.Router, service PersonService, logger *log.Logger) {
	handler := &PersonHandler{
		service: service,
		logger:  logger,
	}

	router.With(render.SetContentType(render.ContentTypeJSON)).Route(personEndpoint, func(r chi.Router) {
		r.Get("/", handler.HandleFindManyPerson)
		r.Post("/", handler.HandleCreatePerson)
		r.Route(personIDEndpoint, func(r chi.Router) {
			r.Get("/", handler.HandleFindPerson)
			r.Put("/", handler.HandleUpdatePerson)
			r.Delete("/", handler.HandleDeletePerson)
		})
	})

	logger.Info("PersonHandler registered", nil)
}

func (h *PersonHandler) HandleCreatePerson(w http.ResponseWriter, r *http.Request) {
	response.OK(w, r, render.M{
		"person": dto.PersonView{},
	})
}

func (h *PersonHandler) HandleFindPerson(w http.ResponseWriter, r *http.Request) {
	response.OK(w, r, render.M{
		"person": dto.PersonView{},
	})
}

func (h *PersonHandler) HandleFindManyPerson(w http.ResponseWriter, r *http.Request) {
	var params dto.FindPersonParams

	if err := params.FromRequest(r); err != nil {
		switch {
		case errors.Is(err, schema.MultiError{}), errors.Is(err, validator.ValidationErrors{}):
			h.logger.Error("bad query params", err, nil)

			response.OK(w, r, render.M{"errors": err})
		default:
			h.logger.Error("something bad happened", err, nil)

			response.OK(w, r, render.M{"error": "internal server error"})
		}

		return
	}

	filter := entity.Person{
		Name:        params.Name,
		Surname:     params.Surname,
		Patronymic:  params.Patronymic,
		Gender:      params.Gender,
		Nationality: params.Nationality,
		Age:         params.Age,
	}

	people, err := h.service.FindMany(r.Context(), params.Limit, params.Offset, filter)
	if err != nil {
		h.logger.Error("failed to find people", err, nil)

		response.OK(w, r, render.M{"error": "failed to find people"})

		return
	}

	response.OK(w, r, render.M{
		"people": converter.PersonEntitiesToViews(people),
	})
}

func (h *PersonHandler) HandleUpdatePerson(w http.ResponseWriter, r *http.Request) {
	response.OK(w, r, render.M{
		"person": dto.PersonView{},
	})
}

func (h *PersonHandler) HandleDeletePerson(w http.ResponseWriter, r *http.Request) {
	response.OK(w, r, render.M{
		"person_id": "",
	})
}
