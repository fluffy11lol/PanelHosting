package models

import "errors"

var (
	ErrInvalidPassword   = errors.New("invalid password")
	ErrUserExist         = errors.New("user already exist, try other name")
	ErrUserNotExist      = errors.New("no user with this name")
	ErrInvalidExpression = errors.New("invalid format of expression")
	ErrEmptyField        = errors.New("empty fields not allowed")
)
