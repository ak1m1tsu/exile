package update

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/insan1a/exile/internal/lib/sl"
	"github.com/insan1a/exile/internal/lib/validator"
	"github.com/insan1a/exile/internal/models"
	"github.com/insan1a/exile/internal/server/http/api/response"
	"github.com/insan1a/exile/internal/storage/person"
)

//go:generate go run github.com/vektra/mockery/v2@v2.20.2 --name PersonUpdater --output ./mocks --outpkg mocks
type PersonUpdater interface {
	Update(context.Context, *models.Person) error
}

func New(log *slog.Logger, updater PersonUpdater) func(http.ResponseWriter, *http.Request) {
	type req struct {
		Name        string `json:"name" validate:"omitempty,alpha"`
		Surname     string `json:"surname" validate:"omitempty,alpha"`
		Patronymic  string `json:"patronymic" validate:"omitempty,alpha"`
		Age         int    `json:"age" validate:"omitempty,gte=0,lte=150"`
		Gender      string `json:"gender" validate:"omitempty,oneof=male female"`
		Nationality string `json:"nationality" validate:"omitempty,alpha,len=2"`
	}

	type resp struct {
		response.Response
		Person *models.Person `json:"person,omitempty"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		log := log.With(
			slog.String("request_id", middleware.GetReqID(r.Context())),
			slog.String("person_id", id),
		)

		var input req
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			msg := "invalid request"

			log.Error(msg, sl.Err(err))

			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, resp{Response: response.Error(msg)})

			return
		}

		if err := validator.ValidateStruct(input); err != nil {
			msg := "invalid request"

			log.Error(msg, sl.Err(err))

			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, resp{Response: response.Error(msg)})

			return
		}

		p := models.Person{
			ID:          id,
			Name:        input.Name,
			Surname:     input.Surname,
			Patronymic:  input.Patronymic,
			Age:         input.Age,
			Gender:      input.Gender,
			Nationality: input.Nationality,
		}

		if err := updater.Update(r.Context(), &p); err != nil {
			if errors.Is(err, person.ErrNotFound) {
				msg := "the person not found"

				log.Error(msg, sl.Err(err))

				render.Status(r, http.StatusNotFound)
				render.JSON(w, r, resp{Response: response.Error(msg)})

				return
			}

			msg := "failed to update the person"

			log.Error(msg, sl.Err(err))

			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, resp{Response: response.Error(msg)})

			return
		}

		log.Info("the person updated", slog.Any("person", p))

		render.JSON(w, r, resp{
			Response: response.OK(),
			Person:   &p,
		})
	}
}
