package list

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/gorilla/schema"
	"github.com/insan1a/exile/internal/lib/sl"
	"github.com/insan1a/exile/internal/lib/validator"
	"github.com/insan1a/exile/internal/models"
	"github.com/insan1a/exile/internal/server/http/api/response"
	"github.com/insan1a/exile/internal/storage/person"
)

//go:generate go run github.com/vektra/mockery/v2@v2.20.2 --name PersonLister --output ./mocks --outpkg mocks
type PersonLister interface {
	List(ctx context.Context, filter *models.Filter, query string) ([]models.Person, error)
}

func New(log *slog.Logger, lister PersonLister) func(http.ResponseWriter, *http.Request) {
	type res struct {
		response.Response
		People []models.Person `json:"people,omitempty"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		log := log.With(
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		if err := r.ParseForm(); err != nil {
			msg := "invalid request"

			log.Error(msg, sl.Err(err))

			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, res{Response: response.Error(msg)})

			return
		}

		filter := new(models.Filter)
		if err := schema.NewDecoder().Decode(filter, r.Form); err != nil {
			msg := "invalid request"

			log.Error(msg, sl.Err(err))

			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, res{Response: response.Error(msg)})

			return
		}

		if err := validator.ValidateStruct(*filter); err != nil {
			msg := "invalid request"

			log.Error(msg, sl.Err(err))

			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, res{Response: response.Error(err.Error())})

			return
		}

		p, err := lister.List(r.Context(), filter, r.URL.Query().Encode())
		if err != nil {
			if errors.Is(err, person.ErrNotFoundMany) {
				msg := "people not found"

				log.Error(msg, sl.Err(err), slog.Any("filter", filter))

				render.Status(r, http.StatusNotFound)
				render.JSON(w, r, res{Response: response.Error(msg)})

				return
			}
			msg := "failed to get people"

			log.Error(msg, sl.Err(err), slog.Any("filter", filter))

			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, res{Response: response.Error(msg)})

			return
		}

		log.Info("people found", slog.Any("people", p))

		render.JSON(w, r, res{
			Response: response.OK(),
			People:   p,
		})
	}
}
