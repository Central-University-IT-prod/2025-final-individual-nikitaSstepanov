package validator

import (
	"github.com/go-playground/validator/v10"
	e "github.com/nikitaSstepanov/tools/error"
)

func Struct(s interface{}, args ...Arg) e.Error {
	validate := validator.New()
	/*if err := setupArgs(validate, args); err != nil {
		return e.BadInputErr.WithErr(err)
	}*/

	err := validate.Struct(s)
	if err != nil {
		errors := err.(validator.ValidationErrors)

		return e.BadInputErr.WithErr(errors)
	}

	return nil
}

func StringLength(s string, min int, max int) e.Error {
	if len(s) < min || len(s) > max {
		return lenErr
	}

	return nil
}

func UUID(s string) e.Error {
	toValidate := uuid{
		Value: s,
	}

	validate := validator.New()

	if err := validate.Struct(toValidate); err != nil {
		return uuidErr
	}

	return nil
}
