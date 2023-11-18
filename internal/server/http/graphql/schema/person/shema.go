package person

import (
	"log/slog"

	"github.com/graphql-go/graphql"
	"github.com/insan1a/exile/internal/server/http/handlers/person/delete"
	"github.com/insan1a/exile/internal/server/http/handlers/person/get"
	"github.com/insan1a/exile/internal/server/http/handlers/person/list"
	"github.com/insan1a/exile/internal/server/http/handlers/person/save"
	"github.com/insan1a/exile/internal/server/http/handlers/person/update"
)

//go:generate go run github.com/vektra/mockery/v2@v2.20.2 --name Storage --output ./mocks --outpkg mocks
type PeopleServicer interface {
	get.PersonGetter
	save.PersonSaver
	list.PersonLister
	delete.PersonDeleter
	update.PersonUpdater
}

func New(log *slog.Logger, svc PeopleServicer) (graphql.Schema, error) {
	personType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Person",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.String,
			},
			"name": &graphql.Field{
				Type: graphql.String,
			},
			"surname": &graphql.Field{
				Type: graphql.String,
			},
			"patronymic": &graphql.Field{
				Type: graphql.String,
			},
			"age": &graphql.Field{
				Type: graphql.Int,
			},
			"gender": &graphql.Field{
				Type: graphql.String,
			},
			"nationality": &graphql.Field{
				Type: graphql.String,
			},
		},
	})

	mutation := graphql.NewObject(graphql.ObjectConfig{
		Name: "RootMutation",
		Fields: graphql.Fields{
			"createPerson": &graphql.Field{
				Type:        personType,
				Description: "Create person",
				Args: graphql.FieldConfigArgument{
					"name": &graphql.ArgumentConfig{
						Type:        graphql.NewNonNull(graphql.String),
						Description: "Name",
					},
					"surname": &graphql.ArgumentConfig{
						Type:        graphql.NewNonNull(graphql.String),
						Description: "Surname",
					},
					"patronymic": &graphql.ArgumentConfig{
						Type:        graphql.String,
						Description: "Patronymic",
					},
				},
				Resolve: Save(log, svc),
			},
			"updatePerson": &graphql.Field{
				Type:        personType,
				Description: "Update person",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type:        graphql.NewNonNull(graphql.String),
						Description: "Id",
					},
					"name": &graphql.ArgumentConfig{
						Type:        graphql.String,
						Description: "Name",
					},
					"surname": &graphql.ArgumentConfig{
						Type:        graphql.String,
						Description: "Surname",
					},
					"patronymic": &graphql.ArgumentConfig{
						Type:        graphql.String,
						Description: "Patronymic",
					},
					"age": &graphql.ArgumentConfig{
						Type:        graphql.Int,
						Description: "Age",
					},
					"gender": &graphql.ArgumentConfig{
						Type:        graphql.String,
						Description: "Gender",
					},
					"nationality": &graphql.ArgumentConfig{
						Type:        graphql.String,
						Description: "Nationality",
					},
				},
				Resolve: Update(log, svc),
			},
			"deletePerson": &graphql.Field{
				Type:        personType,
				Description: "Delete person",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type:        graphql.NewNonNull(graphql.String),
						Description: "Id",
					},
				},
				Resolve: Delete(log, svc),
			},
		},
	})

	query := graphql.NewObject(graphql.ObjectConfig{
		Name: "RootQuery",
		Fields: graphql.Fields{
			"person": &graphql.Field{
				Type:        personType,
				Description: "Get person",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type:        graphql.NewNonNull(graphql.String),
						Description: "Id",
					},
				},
				Resolve: One(log, svc),
			},
			"people": &graphql.Field{
				Type:        graphql.NewList(personType),
				Description: "Get people",
				Args: graphql.FieldConfigArgument{
					"name": &graphql.ArgumentConfig{
						Type:        graphql.String,
						Description: "Name",
					},
					"surname": &graphql.ArgumentConfig{
						Type:        graphql.String,
						Description: "Surname",
					},
					"patronymic": &graphql.ArgumentConfig{
						Type:        graphql.String,
						Description: "Patronymic",
					},
					"age": &graphql.ArgumentConfig{
						Type:        graphql.Int,
						Description: "Age",
					},
					"gender": &graphql.ArgumentConfig{
						Type:        graphql.String,
						Description: "Gender",
					},
					"nationality": &graphql.ArgumentConfig{
						Type:        graphql.String,
						Description: "Nationality",
					},
					"limit": &graphql.ArgumentConfig{
						Type:        graphql.Int,
						Description: "Limit",
					},
					"skip": &graphql.ArgumentConfig{
						Type:        graphql.Int,
						Description: "Skip",
					},
				},
				Resolve: List(log, svc),
			},
		},
	})

	return graphql.NewSchema(graphql.SchemaConfig{
		Query:    query,
		Mutation: mutation,
	})
}
