package utils

import (
	"errors"
	"fmt"
)

var (
	ErrUniqueLink  = errors.New(`this link already exists`)
	ErrDeletedLink = errors.New(`this link has been removed`)
)

type (
	InsertUniqueLinkError struct {
		Link string
		Err  error
	}

	DeletedLinkError struct {
		ShortURL string
		Err      error
	}
)

func NewInsertUniqueLinkError(l string) error {
	return &InsertUniqueLinkError{
		Err:  ErrUniqueLink,
		Link: l,
	}
}

func NewDeletedLinkError(su string) error {
	return &DeletedLinkError{
		Err:      ErrDeletedLink,
		ShortURL: su,
	}
}

func (iu *InsertUniqueLinkError) Error() string {
	return fmt.Sprintf("%v: %v", iu.Err, iu.Link)
}

func (de *DeletedLinkError) Error() string {
	return fmt.Sprintf("%v: %v", de.Err, de.ShortURL)
}

func (iu *InsertUniqueLinkError) Unwrap() error {
	return iu.Err
}

func (de *DeletedLinkError) Unwrap() error {
	return de.Err
}
