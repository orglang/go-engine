package core

import (
	"errors"
	"fmt"

	"orglang/go-engine/internal/adt/identity"
)

var (
	ErrTermDoesNotExist               = errors.New("process term does not exist")
	ErrDefinitionMissing              = errors.New("process term has no Definition")
	ErrTermExpRepositoryNotConfigured = errors.New("term Definition repository is not configured")
)

func ErrTermMissing(termID identity.ADT) error {
	return fmt.Errorf("%w: %v", ErrTermDoesNotExist, termID)
}

func DefinitionMissing(termID identity.ADT) error {
	return fmt.Errorf("%w: %v", ErrDefinitionMissing, termID)
}
