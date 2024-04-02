package configutils

import (
	"fmt"
)

type ValueNotAllowedError struct {
	name     string
	received string
	allowed  []string
}

func (e ValueNotAllowedError) Error() string {
	return fmt.Sprintf("Invalid value for %s; Expected one of %v; Got %s", e.name, e.allowed, e.received)
}

func NewErrValueNotAllowed(varName, received string, allowed []string) ValueNotAllowedError {
	return ValueNotAllowedError{
		name:     varName,
		received: received,
		allowed:  allowed,
	}
}
