package internal

import "errors"

var (
	ErrPepito error = errors.New("pepito")
	ErrInput  error = errors.New("missing input values")
)
