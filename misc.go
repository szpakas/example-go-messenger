package main

import "strings"

// NewValidationError constructs new validation error from different elements (also from other ValidationErrors)
func NewValidationError(errs ...interface{}) *ValidationError {
	e := ValidationError{}
	for _, eIn := range errs {
		switch eInCst := eIn.(type) {
		case string:
			e.msgs = append(e.msgs, eInCst)
		case *ValidationError:
			e.msgs = append(e.msgs, eInCst.msgs...)
		case ValidationError:
			e.msgs = append(e.msgs, eInCst.msgs...)
		}
	}
	return &e
}

// ValidationError is an Error triggered by validation failure
type ValidationError struct {
	msgs []string
}

// Error returns human readable representation of error
func (e ValidationError) Error() string {
	return strings.Join(append([]string{"validation failed"}, e.msgs...), ": ")
}
