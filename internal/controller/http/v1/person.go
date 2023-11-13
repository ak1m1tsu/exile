package v1

import (
	"context"
	"errors"
	"net/http"

	chi "github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/schema"
	"github.com/insan1a/exile/internal/controller/http/dto"
	"github.com/insan1a/exile/internal/controller/http/response"
	"github.com/insan1a/exile/internal/converter"
	"github.com/insan1a/exile/internal/domain/entity"
	"github.com/insan1a/exile/internal/domain/service"
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

type personHandlerFunc func(w http.ResponseWriter, r *http.Request, personID string)

type personHandler struct {
	service PersonService
	logger  *log.Logger
}

func MountPersonHandler(router chi.Router, service PersonService, logger *log.Logger) {
	handler := &personHandler{
		service: service,
		logger:  logger,
	}

	router.With(render.SetContentType(render.ContentTypeJSON)).Route(personEndpoint, func(r chi.Router) {
		r.Get("/", handler.handleFindManyPerson)
		r.Post("/", handler.handleCreatePerson)
		r.Route(personIDEndpoint, func(r chi.Router) {
			r.Get("/", handler.personContext(handler.handleFindPerson))
			r.Put("/", handler.personContext(handler.handleUpdatePerson))
			r.Delete("/", handler.personContext(handler.HandleDeletePerson))
		})
	})

	logger.Info("personHandler registered", nil)
}

func (h *personHandler) personContext(next personHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		personID := chi.URLParam(r, personIDParam)

		next(w, r, personID)
	}
}

func (h *personHandler) handleFindManyPerson(w http.ResponseWriter, r *http.Request) {
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

	filter := converter.FindPersonParamsToEntity(params)

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

func (h *personHandler) handleFindPerson(w http.ResponseWriter, r *http.Request, personID string) {
	person, err := h.service.FindByID(r.Context(), personID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrPersonNotFound):
			h.logger.Error("person not found", err, nil)

			response.NotFound(w, r)
		default:
			h.logger.Error("something bad happened", err, nil)

			response.InternalServerError(w, r)
		}

		return
	}

	response.OK(w, r, render.M{
		"person": converter.PersonEntityToView(person),
	})
}

func (h *personHandler) handleCreatePerson(w http.ResponseWriter, r *http.Request) {
	var body dto.CreatePersonDTO
	if err := body.FromRequest(r); err != nil {
		switch {
		case errors.Is(err, validator.ValidationErrors{}):
			h.logger.Error("bad request body", err, nil)

			response.UnprocessableEntity(w, r, err)
		default:
			h.logger.Error("something bad happened", err, nil)

			response.InternalServerError(w, r)
		}

		return
	}

	person, err := h.service.Store(r.Context(), converter.CreatePersonDTOToEntity(body))
	if err != nil {
		h.logger.Error("failed to store person", err, nil)

		response.InternalServerError(w, r)

		return
	}

	response.OK(w, r, render.M{
		"person": converter.PersonEntityToView(person),
	})
}

func (h *personHandler) handleUpdatePerson(w http.ResponseWriter, r *http.Request, personID string) {
	var body dto.UpdatePersonDTO
	if err := body.FromRequest(r); err != nil {
		switch {
		case errors.Is(err, &validator.ValidationErrors{}):
			h.logger.Error("bad request body", err, nil)

			response.UnprocessableEntity(w, r, err)
		default:
			h.logger.Error("something bad happened", err, nil)

			response.InternalServerError(w, r)
		}

		return
	}

	person, err := h.service.Update(r.Context(), converter.UpdatePersonDTOToEntity(personID, body))
	if err != nil {
		h.logger.Error("failed to update person", err, nil)

		response.InternalServerError(w, r)

		return
	}

	response.OK(w, r, render.M{
		"person": converter.PersonEntityToView(person),
	})
}

func (h *personHandler) HandleDeletePerson(w http.ResponseWriter, r *http.Request, personID string) {
	if err := h.service.Delete(r.Context(), personID); err != nil {
		h.logger.Error("failed to delete person", err, nil)

		response.InternalServerError(w, r)

		return
	}

	response.OK(w, r, render.M{
		"person_id": personID,
	})
}
