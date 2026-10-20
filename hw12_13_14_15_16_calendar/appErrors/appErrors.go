package appErrors

import "errors"

var (
	ErrIdDoesNotExist        = errors.New("id does not exist")
	ErrTitleIsNotSet         = errors.New("title is not set")
	ErrDateStartIsNotSet     = errors.New("date start is not set")
	ErrDateEndIsNotSet       = errors.New("date end is not set")
	ErrDateStartBelowDateEnd = errors.New("date start is after date end")
	ErrUserIdIsNotSet        = errors.New("user id is not set")
	ErrNotifyTimeIsNegative  = errors.New("notification time is negative")
	ErrValidation            = errors.New("validation error")
)
