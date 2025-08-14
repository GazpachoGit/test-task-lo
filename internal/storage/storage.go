package storage

import "errors"

var (
	ErrTaskNotFound      = errors.New("task not found")
	ErrTaskStatusInvalid = errors.New("invalid task status")
)

const IDLength = 6
