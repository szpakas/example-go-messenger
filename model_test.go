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

// -- section: Message
func Test_Model_Msg_Success(t *testing.T) {
	a.NoError(t, tfMsgAA.Validate())
}

func Test_Model_Msg_Failure(t *testing.T) {
	tests := map[string]struct {
		msg  Message
		eStr string
	}{
		"no ID":                  {tfMsgAXA_NoID, "missing ID"},
		"no Body":                {tfMsgAXB_NoBody, "missing Body"},
		"no AuthorID":            {tfMsgXXA_NoAuthorID, "missing AuthorID"},
		"invalid Tag: empty":     {tfMsgAXC_NoTag, "invalid Tag: empty value"},
		"invalid Tag: too short": {tfMsgAXD_TagTooShort, "invalid Tag: too short"},
	}

	for s, tc := range tests {
		a.EqualError(t, tc.msg.Validate(), fmt.Sprintf("validation failed: %s", tc.eStr), "case: %s", s)
	}
}
