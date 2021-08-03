package seng

import "errors"

var ErrNotFoundRoute = errors.New("not found Path")
var ErrQueryParamNotFound = errors.New("query param not found")
