package errors_rep

import "errors"

var (
	ErrNoRowsAffected = errors.New("No rows were affected")
	ErrNullValue      = errors.New("Null value")
)
