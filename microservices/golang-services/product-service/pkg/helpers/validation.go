package helpers

import (
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func InitValidator() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		registerCustomValidators(v)
	}
}

func registerCustomValidators(v *validator.Validate) {
	v.RegisterValidation("not_blank", func(fl validator.FieldLevel) bool {
		val := fl.Field().String()
		return strings.TrimSpace(val) != ""
	})
}
func ParseValidationErrors(err error) map[string][]string {
	res := map[string][]string{}

	if errs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range errs {
			field := e.Field()
			msg := validationMessage(e)
			res[field] = append(res[field], msg)
		}
	} else {
		res["error"] = []string{err.Error()}
	}

	return res
}

func validationMessage(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return "The " + e.Field() + " field is required."
	case "email":
		return "The " + e.Field() + " format is invalid."
	case "min":
		return "The " + e.Field() + " must be at least " + e.Param() + " characters."
	case "max":
		return "The " + e.Field() + " must be at most " + e.Param() + " characters."
	case "gte":
		return "The " + e.Field() + " must be greater than or equal to " + e.Param() + "."
	default:
		return "The " + e.Field() + " is invalid."
	}
}
