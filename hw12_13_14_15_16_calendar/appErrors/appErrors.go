package appErrors

import "errors"

var ErrIdDoesNotExist = errors.New("id does not exist")
var ErrTitleIsNotSet = errors.New("title is not set")
var ErrDateStartIsNotSet = errors.New("date start is not set")
var ErrDateEndIsNotSet = errors.New("date end is not set")
var ErrDateStartBelowDateEnd = errors.New("date start below date end")
var ErrUserIdIsNotSet = errors.New("user id is not set")
