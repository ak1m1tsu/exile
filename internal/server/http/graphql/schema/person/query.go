package person

import (
	"context"
	"log/slog"

	"github.com/graphql-go/graphql"
	"github.com/insan1a/exile/internal/lib/sl"
	"github.com/insan1a/exile/internal/lib/validator"
	"github.com/insan1a/exile/internal/models"
	"github.com/insan1a/exile/internal/server/http/handlers/person/get"
	"github.com/insan1a/exile/internal/server/http/handlers/person/list"
	"github.com/mitchellh/mapstructure"
)

func List(log *slog.Logger, listter list.PersonLister) func(graphql.ResolveParams) (interface{}, error) {
	return func(params graphql.ResolveParams) (interface{}, error) {
		var input models.Filter
		if err := mapstructure.Decode(params.Args, &input); err != nil {
			msg := "invalid request"

			log.Error(msg, sl.Err(err), slog.Any("args", params.Args))

			return nil, err
		}

		if err := validator.ValidateStruct(input); err != nil {
			msg := "failed to validate request"

			log.Error(msg, sl.Err(err), slog.Any("input", input), slog.Any("args", params.Args))

			return nil, err
		}

		p, err := listter.List(context.Background(), &input, input.String())
		if err != nil {
			msg := "failed to find people"

			log.Error(msg, sl.Err(err), slog.Any("input", input), slog.Any("args", params.Args))

			return nil, err
		}

		log.Info("people found", slog.Any("input", input), slog.Any("args", params.Args))

		return p, nil
	}
}

func One(log *slog.Logger, getter get.PersonGetter) func(graphql.ResolveParams) (interface{}, error) {
	type req struct {
		ID string `mapstructure:"id" validate:"required,uuid"`
	}
	return func(params graphql.ResolveParams) (interface{}, error) {
		var input req
		if err := mapstructure.Decode(params.Args, &input); err != nil {
			msg := "invalid request"

			log.Error(msg, sl.Err(err), slog.Any("args", params.Args))

			return nil, err
		}

		if err := validator.ValidateStruct(input); err != nil {
			msg := "failed to validate request"

			log.Error(msg, sl.Err(err), slog.Any("input", input), slog.Any("args", params.Args))

			return nil, err
		}

		p, err := getter.Get(context.Background(), input.ID)
		if err != nil {
			msg := "failed to get person"

			log.Error(msg, sl.Err(err), slog.Any("input", input), slog.Any("args", params.Args))

			return nil, err
		}

		log.Info("person found", slog.Any("input", input), slog.Any("args", params.Args), slog.Any("person", p))

		return p, nil
	}
}
