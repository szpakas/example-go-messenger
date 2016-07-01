package main

import (
	"testing"

	"fmt"
	a "github.com/stretchr/testify/assert"
	ar "github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strings"

	"encoding/json"
	"github.com/davecgh/go-spew/spew"
)

var _ = spew.Config

func Test_HTTPServer_Factory(t *testing.T) {
	st := NewMemoryStorage()
	h := NewHTTPDefaultHandler(st)
	s := NewHTTPServer("A12345.example.com", 9876, h)

	ar.NotNil(t, s, "empty element returned")

	a.Equal(t, "A12345.example.com:9876", s.Addr, "server address mismatch")

	//ar.IsType(t, &httpHandler{}, s.Handler, "mismatch in handler type")
	//hCst := s.Handler.(*httpHandler)
	//a.Equal(t, st, hCst.Storer, "storage not attached to handler")
}

func Test_HTTPHandler_User_Create_Success(t *testing.T) {
	st := NewMemoryStorage()
	h := NewHTTPDefaultHandler(st)
	ts := httptest.NewServer(h)
	defer ts.Close()

	bR := strings.NewReader(tfTrInUserA_JSON)
	res, err := http.Post(fmt.Sprintf("%s/v1/users", ts.URL), "application/json", bR)

	// THEN: validate response
	ar.NoError(t, err, "unexpected error from HTTP client")
	a.EqualValues(t, 0, res.ContentLength, "non empty response body")
	a.Equal(t, http.StatusCreated, res.StatusCode, "mismatch on response code")

	// AND: validate storage
	userGot, err := st.UserFindByName(tfTrInUserA.Name)
	ar.NoError(t, err, "unexpected error on user seek")
	a.Equal(t, tfTrInUserA.Name, userGot.Name, "User: mismatch in Name")
	a.NotZero(t, userGot.ID, "User: zero ID")
}

func Test_HTTPHandler_Message_Create_Success(t *testing.T) {
	st := NewMemoryStorage()
	h := NewHTTPDefaultHandler(st)
	ts := httptest.NewServer(h)
	defer ts.Close()

	// GIVEN: author match existing user
	ar.NoError(t, st.UserSave(&tfUserA))

	bR := strings.NewReader(tfTrInMsgAA_JSON)
	res, err := http.Post(fmt.Sprintf("%s/v1/messages", ts.URL), "application/json", bR)

	// THEN: validate response
	ar.NoError(t, err, "unexpected error from HTTP client")
	a.EqualValues(t, 0, res.ContentLength, "non empty response body")
	a.Equal(t, http.StatusCreated, res.StatusCode, "mismatch on response code")

	// AND: validate tag association
	msgsIDs, err := st.MsgsIDsFindByTag(tfTrInMsgAA.Tag)
	ar.NoError(t, err, "unexpected error on tag seek")
	ar.Len(t, msgsIDs, 1, "incorrect number of messages associated to tag returned")

	// AND: validate message stored
	msgGot, err := st.MsgLoad(msgsIDs[0])
	ar.NoError(t, err, "unexpected error on message load")
	a.Equal(t, tfMsgAA.AuthorID, msgGot.AuthorID, "message AuthorID mismatch")
	a.Equal(t, tfMsgAA.Body, msgGot.Body, "message Body mismatch")
	a.Equal(t, tfMsgAA.Tag, msgGot.Tag, "message Tag mismatch")
}

func Test_HTTPHandler_Message_Find_Success_Found(t *testing.T) {
	var ts *httptest.Server
	defer func() {
		if ts != nil {
			ts.Close()
		}
	}()

	tests := map[string]struct {
		dbUsers []User
		dbMsg   []Message
		tag     Tag
		exp     TrOutMessagesCollection
	}{
		"single": {
			[]User{tfUserA},
			[]Message{tfMsgAA},
			tfTrInMsgAA.Tag,
			TrOutMessagesCollection{tfTrOutMsgAA},
		},
		"two": {
			[]User{tfUserA},
			[]Message{tfMsgAA, tfMsgAB},
			tfTrInMsgAA.Tag,
			TrOutMessagesCollection{tfTrOutMsgAA, tfTrOutMsgAB},
		},
		"three, cross user, partial": {
			[]User{tfUserA, tfUserB},
			[]Message{tfMsgAA, tfMsgAB, tfMsgBA, tfMsgBB},
			tfTrInMsgAA.Tag,
			TrOutMessagesCollection{tfTrOutMsgAA, tfTrOutMsgAB, tfTrOutMsgBA},
		},
	}

	for sym, tc := range tests {
		st := NewMemoryStorage()
		h := NewHTTPDefaultHandler(st)
		ts = httptest.NewServer(h)

		// GIVEN: user and message are in DB
		for _, u := range tc.dbUsers {
			uC := u
			ar.NoError(t, st.UserSave(&uC), "case: %s", sym)
		}
		for _, m := range tc.dbMsg {
			mC := m
			ar.NoError(t, st.MsgSave(&mC), "case: %s", sym)
		}

		res, err := http.Get(fmt.Sprintf("%s/v1/messages?tag=%s", ts.URL, tc.tag))

		// THEN: validate response
		if !a.NoError(t, err, "unexpected error from HTTP client") {
			ts.Close()
			continue
		}
		a.Equal(t, http.StatusOK, res.StatusCode, "mismatch on response code")
		a.Equal(t, "application/json", res.Header.Get("Content-Encoding"), "mismatch on response content encoding")

		var resBodyGot TrOutMessagesCollection
		err = json.NewDecoder(res.Body).Decode(&resBodyGot)
		res.Body.Close()
		if !a.NoError(t, err, "unexpected error on response body read") {
			ts.Close()
			continue
		}

		// validate message collection
		if !a.Len(t, resBodyGot, len(tc.exp), "mismatch on number of elements returned") {
			ts.Close()
			continue
		}

		for _, resMsgExp := range tc.exp {
			a.Contains(t, resBodyGot, resMsgExp)
		}

		ts.Close()
	}
}

func Test_HTTPHandler_Message_Find_Success_NotFound(t *testing.T) {
	st := NewMemoryStorage()
	h := NewHTTPDefaultHandler(st)
	ts := httptest.NewServer(h)
	defer ts.Close()

	// GIVEN: no matching messages are in DB
	res, err := http.Get(fmt.Sprintf("%s/v1/messages?tag=%s", ts.URL, "NON_EXISTING_TAG"))

	// THEN: validate response
	ar.NoError(t, err, "unexpected error from HTTP client")
	a.Equal(t, http.StatusNotFound, res.StatusCode, "mismatch on response code")
	a.EqualValues(t, 0, res.ContentLength, "non empty response body")
	a.Zero(t, res.Header.Get("Content-Encoding"), "unexpected content encoding response header")
}
