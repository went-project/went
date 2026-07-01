package http

func ValidatorTemplate() string {
	return `package requests

import "github.com/go-playground/validator/v10"

var validate = validator.New()
`
}
