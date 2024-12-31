package wordports

import "errors"

var (
	ErrWordNotFound      = errors.New("word not found")
	ErrWordAlreadyExists = errors.New("word already exists")
)
