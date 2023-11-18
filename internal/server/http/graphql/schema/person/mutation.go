package person

import (
	"context"
	"log/slog"

	"github.com/graphql-go/graphql"
	"github.com/insan1a/exile/internal/lib/sl"
	"github.com/insan1a/exile/internal/lib/validator"
	"github.com/insan1a/exile/internal/models"
	"github.com/insan1a/exile/internal/server/http/api/response"
	"github.com/insan1a/exile/internal/server/http/handlers/person/delete"
	"github.com/insan1a/exile/internal/server/http/handlers/person/save"
	"github.com/insan1a/exile/internal/server/http/handlers/person/update"
	"github.com/mitchellh/mapstructure"
)

func Save(log *slog.Logger, saver save.PersonSaver) func(params graphql.ResolveParams) (interface{}, error) {
	type req struct {
		Name       string `mapstructure:"name" validate:"required,alpha"`
		Surname    string `mapstructure:"surname" validate:"required,alpha"`
		Patronymic string `mapstructure:"patronymic" validate:"omitempty,alpha"`
	}
	return func(params graphql.ResolveParams) (interface{}, error) {
		var input req
		if err := mapstructure.Decode(params.Args, &input); err != nil {
			msg := "invalid request"

			log.Error(msg, sl.Err(err), slog.Any("args", params.Args))

			return nil, err
		}

		if err := validator.ValidateStruct(input); err != nil {
			msg := "failed to validate request params"

			log.Error(msg, sl.Err(err), slog.Any("input", input))

			return nil, err
		}

		p := models.Person{
			Name:       input.Name,
			Surname:    input.Surname,
			Patronymic: input.Patronymic,
		}
		if err := saver.Save(context.Background(), p); err != nil {
			msg := "failed to save person"

			log.Error(msg, sl.Err(err), slog.Any("input", input))

			return nil, err
		}

		log.Info("the person successfully saved", slog.Any("person", p), slog.Any("input", input))

		return p, nil
	}
}

func Update(log *slog.Logger, updater update.PersonUpdater) func(params graphql.ResolveParams) (interface{}, error) {
	type req struct {
		ID          string `mapstructure:"id" validate:"required,uuid"`
		Name        string `mapstructure:"name" validate:"omitempty,alpha"`
		Surname     string `mapstructure:"surname" validate:"omitempty,alpha"`
		Patronymic  string `mapstructure:"patronymic" validate:"omitempty,alpha"`
		Age         int    `mapstructure:"age" validate:"omitempty,min=0"`
		Gender      string `mapstructure:"gender" validate:"omitempty,alpha"`
		Nationality string `mapstructure:"nationality" validate:"omitempty,alpha"`
	}
	return func(params graphql.ResolveParams) (interface{}, error) {
		var input req
		if err := mapstructure.Decode(params.Args, &input); err != nil {
			msg := "invalid request"

			log.Error(msg, sl.Err(err), slog.Any("args", params.Args))
		}

		if err := validator.ValidateStruct(input); err != nil {
			msg := "failed to validate request params"

			log.Error(msg, sl.Err(err), slog.Any("input", input))

			return nil, err
		}

		p := models.Person{
			ID:          input.ID,
			Name:        input.Name,
			Surname:     input.Surname,
			Patronymic:  input.Patronymic,
			Age:         input.Age,
			Gender:      input.Gender,
			Nationality: input.Nationality,
		}

		if err := updater.Update(context.Background(), &p); err != nil {
			msg := "failed to update person"

			log.Error(msg, sl.Err(err), slog.Any("input", input))

			return nil, err
		}

		log.Info("the person successfully updated", slog.Any("person", p), slog.Any("input", input))

		return p, nil
	}
}

func Delete(log *slog.Logger, deleter delete.PersonDeleter) func(params graphql.ResolveParams) (interface{}, error) {
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
			msg := "failed to validate request params"

			log.Error(msg, sl.Err(err), slog.Any("input", input))

			return nil, err
		}

		if err := deleter.Delete(context.Background(), input.ID); err != nil {
			msg := "failed to delete person"

			log.Error(msg, sl.Err(err), slog.Any("input", input))

			return nil, err
		}

		log.Info("the person successfully deleted", slog.Any("input", input))

		return response.OK(), nil
	}
}
