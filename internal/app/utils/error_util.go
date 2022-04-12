package utils

import (
	"errors"
	"fmt"
)

var ErrUniqueLink = errors.New(`this link already exists`)

type InsertUniqueLinkError struct {
	Link string
	Err  error
}

func (iu *InsertUniqueLinkError) Error() string {
	return fmt.Sprintf("%v: %v", iu.Err, iu.Link)
}

func NewInsertUniqueLinkError(l string) error {
	return &InsertUniqueLinkError{
		Err:  ErrUniqueLink,
		Link: l,
	}
}

func (iu *InsertUniqueLinkError) Unwrap() error {
	return iu.Err
}
