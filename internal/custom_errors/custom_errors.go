package custom_errors

import (
	"errors"
	"fmt"
)

// define error to provide errors.Is()
var ErrKeyNotExists = errors.New("ключ не существует")

// KeyNotExistsError defines custom error
type KeyNotExistsError struct {
	Key string
}

func (e *KeyNotExistsError) Error() string {
	return fmt.Sprintf("ключ '%s' не существует", e.Key)
}

// Unwrap provide to use errors.Is() method
func (e *KeyNotExistsError) Unwrap() error {
	return ErrKeyNotExists
}

// NewKeyNotExistsError creates custom error instance
func NewKeyNotExistsError(key string) error {
	return &KeyNotExistsError{Key: key}
}
