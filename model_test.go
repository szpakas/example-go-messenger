package main

import (
	"fmt"
	"testing"

	a "github.com/stretchr/testify/assert"
)

// -- section: Tag
func Test_Model_Tag_Success(t *testing.T) {
	a.NoError(t, tfTagA.Validate())
}

func Test_Model_Tag_Failure(t *testing.T) {
	tests := map[string]struct {
		tag  Tag
		eStr string
	}{
		"zero":      {Tag(""), "empty value"},
		"too short": {Tag("a"), "too short"},
		"too long": {Tag(func() string {
			s := ""
			for i := 0; i < 256; i++ {
				s += "A"
			}
			return s
		}()), "too long"},
	}

	for s, tc := range tests {
		a.EqualError(t, tc.tag.Validate(), fmt.Sprintf("validation failed: %s", tc.eStr), "case: %s", s)
	}
}
