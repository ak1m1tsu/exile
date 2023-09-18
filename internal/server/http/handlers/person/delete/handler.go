package delete

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/romankravchuk/effective-mobile-test-task/internal/lib/sl"
	"github.com/romankravchuk/effective-mobile-test-task/internal/server/http/api/response"
	"github.com/romankravchuk/effective-mobile-test-task/internal/storage/person"
)

//go:generate go run github.com/vektra/mockery/v2@v2.20.2 --name PersonDeleter --output ./mocks --outpkg mocks
type PersonDeleter interface {
	Delete(ctx context.Context, id string) error
}

func New(log *slog.Logger, deleter PersonDeleter) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		log := log.With(
			slog.String("request_id", middleware.GetReqID(r.Context())),
			slog.String("person_id", id),
		)

		err := deleter.Delete(r.Context(), id)
		if err != nil {
			if errors.Is(err, person.ErrNotFound) {
				msg := "the person not found"

				log.Error(msg, sl.Err(err))

				render.Status(r, http.StatusNotFound)
				render.JSON(w, r, response.Error(msg))

				return
			}

			msg := "failed to delete the person"

			log.Error(msg, sl.Err(err))

			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error(msg))

			return
		}

		render.JSON(w, r, response.OK())
	}
}
