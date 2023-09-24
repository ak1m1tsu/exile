package get

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/romankravchuk/effective-mobile-test-task/internal/lib/sl"
	"github.com/romankravchuk/effective-mobile-test-task/internal/models"
	"github.com/romankravchuk/effective-mobile-test-task/internal/server/http/api/response"
	"github.com/romankravchuk/effective-mobile-test-task/internal/storage/person"
)

//go:generate go run github.com/vektra/mockery/v2@v2.20.2 --name PersonGetter --output ./mocks --outpkg mocks
type PersonGetter interface {
	Get(ctx context.Context, id string) (*models.Person, error)
}

func New(log *slog.Logger, getter PersonGetter) func(http.ResponseWriter, *http.Request) {
	type resp struct {
		response.Response
		*models.Person `json:"person,omitempty"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		log := log.With(
			slog.String("request_id", middleware.GetReqID(r.Context())),
			slog.String("person_id", id),
		)

		p, err := getter.Get(r.Context(), id)
		if err != nil {
			if errors.Is(err, person.ErrNotFound) {
				msg := "the person not found"

				log.Error(msg, sl.Err(err))

				render.Status(r, http.StatusNotFound)
				render.JSON(w, r, resp{Response: response.Error(msg)})

				return
			}

			msg := "failed to get person"

			log.Error(msg, sl.Err(err))

			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, resp{Response: response.Error(msg)})

			return
		}

		log.Info("person found", slog.Any("person", p))

		render.JSON(w, r, resp{
			Response: response.OK(),
			Person:   p,
		})
	}
}
