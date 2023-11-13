package response

import (
	"errors"
	"github.com/go-chi/render"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/schema"
	"net/http"
	"strings"
)

func jsonResponse(w http.ResponseWriter, r *http.Request, statusCode int, data render.M) {
	render.Status(r, statusCode)
	render.JSON(w, r, data)
}

func OK(w http.ResponseWriter, r *http.Request, data render.M) {
	jsonResponse(w, r, http.StatusOK, data)
}

func NotFound(w http.ResponseWriter, r *http.Request) {
	jsonResponse(w, r, http.StatusNotFound, render.M{
		"message": "the resource not found",
	})
}

func MethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	jsonResponse(w, r, http.StatusMethodNotAllowed, render.M{
		"message": "method not allowed",
	})
}

func InternalServerError(w http.ResponseWriter, r *http.Request) {
	jsonResponse(w, r, http.StatusInternalServerError, render.M{
		"message": "internal server error",
	})
}

func BadRequest(w http.ResponseWriter, r *http.Request) {
	jsonResponse(w, r, http.StatusBadRequest, render.M{
		"message": "bad request",
	})
}

func UnprocessableEntity(w http.ResponseWriter, r *http.Request, err error) {
	var (
		data           render.M
		validationErrs *validator.ValidationErrors
		multiError     *schema.MultiError
	)

	switch {
	case errors.As(err, validationErrs):
		errs := make(map[string]string)
		translator := en.New()
		universalTranslator := ut.New(translator, translator)
		trans, _ := universalTranslator.GetTranslator("en")

		for _, e := range *validationErrs {
			errs[strings.ToLower(e.Field())] = strings.ToLower(e.Translate(trans))
		}
		data = render.M{
			"message": "validation error",
			"errors":  errs,
		}
	case errors.As(err, multiError):
		data = render.M{
			"message": "validation error",
			"errors":  multiError,
		}
	}

	jsonResponse(w, r, http.StatusUnprocessableEntity, data)
}
