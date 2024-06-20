package entity

import (
	"errors"
	"fmt"
)

func ErrAlreadyExists(entity string, ID int) error {
	return fmt.Errorf("%s already exists with ID %d", entity, ID)
}
func ErrNotFound(entity string, ID int) error {
	return fmt.Errorf("%s not found with ID %d", entity, ID)
}

var (
	ErrPostCommentsDisable = errors.New("comments on posts are disabled")
	ErrBigContent          = errors.New("the length of the content is more than 2000 characters")
)
