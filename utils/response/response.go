package response

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

type Message struct {
	Message string `json:"message"`
	Detail  any    `json:",omitempty"`
}

var (
	InternalErrResponse = Message{
		Message: "Internal server error",
		Detail:  "An unexpected error occurred while processing your request",
	}
	InputErrResponse = Message{
		Message: "Input error",
		Detail:  nil,
	}
)

func InputFieldError(err error) Message {
	if ok := errors.As(err, &validator.ValidationErrors{}); !ok {
		return InputErrResponse
	}

	var errs []string
	for _, err := range err.(validator.ValidationErrors) {
		errs = append(errs, err.Field()+" is invalid ("+err.Tag()+")")
	}
	return Message{
		Message: "Input error",
		Detail:  errs,
	}
}
