package flam

import (
	"errors"
	"fmt"
)

var (
	ErrNilReference = errors.New("nil reference")

	ErrBagInvalidPath = errors.New("invalid bag path")

	ErrUnknownResource       = errors.New("unknown resource")
	ErrInvalidResourceConfig = errors.New("invalid resource config")
	ErrDuplicateResource     = errors.New("duplicate resource")

	ErrDuplicateProvider = errors.New("duplicate provider")
)

func newErrNilReference(
	arg string,
) error {
	return NewErrorFrom(
		ErrNilReference,
		arg)
}

func newErrBagInvalidPath(
	path string,
) error {
	return NewErrorFrom(
		ErrBagInvalidPath,
		path)
}

func newErrUnknownResource(
	resource string,
	id string,
) error {
	return NewErrorFrom(
		ErrUnknownResource,
		fmt.Sprintf("%s(%s)", resource, id))
}

func newErrInvalidResourceConfig(
	resource string,
	id string,
	config Bag,
) error {
	return NewErrorFrom(
		ErrInvalidResourceConfig,
		fmt.Sprintf("%s(%s) <= %v", resource, id, config))
}

func newErrDuplicateResource(
	id string,
) error {
	return NewErrorFrom(
		ErrDuplicateResource,
		id)
}

func newErrDuplicateProvider(
	id string,
) error {
	return NewErrorFrom(
		ErrDuplicateProvider,
		id)
}
