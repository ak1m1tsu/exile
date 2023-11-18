package person

import (
	"log/slog"
	"net/http"

	"github.com/graphql-go/handler"
	"github.com/insan1a/exile/internal/server/http/graphql/schema/person"
)

func New(log *slog.Logger, svc person.PeopleServicer) http.Handler {
	s, _ := person.New(log, svc)
	return handler.New(&handler.Config{
		Schema:   &s,
		Pretty:   true,
		GraphiQL: false,
	})
}
