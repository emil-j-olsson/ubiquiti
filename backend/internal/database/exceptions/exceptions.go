package exceptions

import "errors"

var (
	ErrorInternal         = errors.New("internal error")
	ErrorNotFound         = errors.New("not found")
	ErrorUnspecifiedField = errors.New("unspecified field")
	ErrorInvalidType      = errors.New("invalid type")
)
