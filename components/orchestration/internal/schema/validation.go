package schema

import (
	"reflect"
	"strings"

	"github.com/formancehq/orchestration/internal/workflow/stages"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func ValidateRequirements(schema stages.Stage) error {
	return validate.Struct(schema)
}

func reportOneOfError(sl validator.StructLevel, value reflect.Value) {
	choices := make([]string, 0)
	for i := 0; i < value.Type().NumField(); i++ {
		jsonTag := value.Type().Field(i).Tag.Get("json")
		choices = append(choices, strings.Split(jsonTag, ",")[0])
	}
	for i := 0; i < value.Type().NumField(); i++ {
		sl.ReportError(
			value.Field(i).Interface(),
			strings.Split(value.Type().Field(i).Tag.Get("json"), ",")[0],
			value.Field(i).Type().Name(),
			strings.Join(choices, " or "),
			"",
		)
	}
}

func oneOf(sl validator.StructLevel) {
	object := sl.Current().Interface()
	valueOfObject := reflect.ValueOf(object)
	defined := false
	for i := 0; i < valueOfObject.Type().NumField(); i++ {
		if !valueOfObject.Field(i).IsZero() {
			if defined {
				reportOneOfError(sl, valueOfObject)
				return
			}
			defined = true
		}
	}
	if !defined {
		reportOneOfError(sl, valueOfObject)
	}
}

func RegisterOneOf(to ...any) {
	for _, v := range to {
		validate.RegisterStructValidation(oneOf, v)
	}
}
