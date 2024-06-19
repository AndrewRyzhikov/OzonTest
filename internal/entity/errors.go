package entity

import (
	"errors"
)

var (
	ErrAlreadyExists = errors.New("entity already exists")
	ErrNotFound      = errors.New("entity not found")
)
