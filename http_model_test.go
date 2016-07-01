package main

import (
	"encoding/json"
	"fmt"
	"testing"

	a "github.com/stretchr/testify/assert"
	ar "github.com/stretchr/testify/require"
)

// -- section: User
func Test_HTTPModel_TrInUser_Validate_Success(t *testing.T) {
	a.NoError(t, tfTrInUserA.Validate())
}

func Test_HTTPModel_TrInUser_Validate_Failure(t *testing.T) {
	tests := map[string]struct {
		obj  TrInUser
		eStr string
	}{
		"no Name":        {TrInUser{}, "Name missing"},
		"name too short": {TrInUser{Name: "A"}, "Name too short"},
	}

	for s, tc := range tests {
		a.EqualError(t, tc.obj.Validate(), fmt.Sprintf("validation failed: %s", tc.eStr), "case: %s", s)
	}
}

func Test_HTTPModel_TrInUser_JSONEncode(t *testing.T) {
	enc, err := json.Marshal(&tfTrInUserA)
	ar.NoError(t, err)
	a.JSONEq(t, tfTrInUserA_JSON, string(enc))
}

// -- section: Message
func Test_HTTPModel_TrInMsg_Validate_Success(t *testing.T) {
	a.NoError(t, tfTrInMsgAA.Validate())
}

func Test_HTTPModel_TrInMsg_Validate_Failure(t *testing.T) {
	tests := map[string]struct {
		obj  TrInMessage
		eStr string
	}{
		"no Body":                {tfTrInMsgAXB_NoBody, "missing Body"},
		"no Author":              {tfTrInMsgXXA_NoAuthor, "missing Author"},
		"invalid Tag: empty":     {tfTrInMsgAXC_NoTag, "invalid Tag: empty value"},
		"invalid Tag: too short": {tfTrInMsgAXD_TagTooShort, "invalid Tag: too short"},
	}

	for s, tc := range tests {
		a.EqualError(t, tc.obj.Validate(), fmt.Sprintf("validation failed: %s", tc.eStr), "case: %s", s)
	}
}

func Test_HTTPModel_TrInMsg_JSONEncode(t *testing.T) {
	enc, err := json.Marshal(&tfTrInMsgAA)
	ar.NoError(t, err)
	a.JSONEq(t, tfTrInMsgAA_JSON, string(enc))
}

func Test_HTTPModel_TrOutMsg_JSONEncode(t *testing.T) {
	enc, err := json.Marshal(&tfTrOutMsgAA)
	ar.NoError(t, err)
	a.JSONEq(t, tfTrOutMsgAA_JSON, string(enc))
}
