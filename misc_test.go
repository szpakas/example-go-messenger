package main

import (
	"testing"

	a "github.com/stretchr/testify/assert"
)

func Test_ValidationError_Simple(t *testing.T) {
	a.EqualError(t, NewValidationError("simpleA"), "validation failed: simpleA")
}

func Test_ValidationError_Chained_Ref_Single(t *testing.T) {
	a.EqualError(t, NewValidationError(NewValidationError("simpleA")), "validation failed: simpleA")
}

func Test_ValidationError_Chained_Ref_Multiple(t *testing.T) {
	a.EqualError(t, NewValidationError(NewValidationError("simpleA"), "simpleB"), "validation failed: simpleA: simpleB")
}

func Test_ValidationError_Chained_Val_Single(t *testing.T) {
	a.EqualError(t, NewValidationError(ValidationError{[]string{"simpleA"}}), "validation failed: simpleA")
}

func Test_ValidationError_Chained_Val_Multiple(t *testing.T) {
	a.EqualError(t, NewValidationError(ValidationError{[]string{"simpleA"}}, "simpleB"), "validation failed: simpleA: simpleB")
}
