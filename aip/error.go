package aip

import "errors"

var (
	ErrInvalidStatus    = errors.New("Invalid status")
	ErrInvalidParameter = errors.New("Invalid parameter")
	ErrInvalidToken     = errors.New("Invalid token")
)
