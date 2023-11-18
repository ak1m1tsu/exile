package save

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/insan1a/exile/internal/lib/sl"
	"github.com/insan1a/exile/internal/lib/validator"
	"github.com/insan1a/exile/internal/models"
	"github.com/insan1a/exile/internal/server/http/api/response"
)

//go:generate go run github.com/vektra/mockery/v2@v2.20.2 --name PersonSaver --output ./mocks --outpkg mocks
type PersonSaver interface {
	Save(ctx context.Context, p models.Person) error
}

func New(log *slog.Logger, saver PersonSaver) func(http.ResponseWriter, *http.Request) {
	type req struct {
		Name       string `json:"name" validate:"required,alpha"`
		Surname    string `json:"surname" validate:"required,alpha"`
		Patronymic string `json:"patronymic" validate:"omitempty,alpha"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		log := log.With(
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var input req
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			msg := "invalid request"

			log.Error(msg, sl.Err(err))

			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error(msg))

			return
		}

		if err := validator.ValidateStruct(input); err != nil {
			msg := "failed to validate request"

			log.Error(msg, sl.Err(err), slog.Any("request_body", input))

			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error(err.Error()))

			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 100*time.Millisecond)
		defer cancel()

		err := saver.Save(ctx, models.Person{
			Name:       input.Name,
			Surname:    input.Surname,
			Patronymic: input.Patronymic,
		})
		if err != nil {
			msg := "failed to save person"

			log.Error(msg, sl.Err(err), slog.Any("request_body", input))

			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error(msg))

			return
		}

		log.Info("the person successfully saved", slog.Any("request_body", input))

		render.JSON(w, r, response.OK())
	}
}
