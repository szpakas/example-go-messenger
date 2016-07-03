package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	a "github.com/stretchr/testify/assert"
	ar "github.com/stretchr/testify/require"
)

func Test_HTTPServer_Factory(t *testing.T) {
	st := NewMemoryStorage()
	h := NewHTTPDefaultHandler(st)
	s := NewHTTPServer("A12345.example.com", 9876, h)

	ar.NotNil(t, s, "empty element returned")

	a.Equal(t, "A12345.example.com:9876", s.Addr, "server address mismatch")
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

func Test_HTTPHandler_User_Create_Failure(t *testing.T) {
	var ts *httptest.Server
	defer func() {
		if ts != nil {
			ts.Close()
		}
	}()

	tests := map[string]struct {
		reqBody     string
		dbUsers     []User
		usErr       error
		usCalledExp bool
		resStatus   int
	}{
		"validation: empty object/name empty": {
			reqBody:   `{}`,
			resStatus: http.StatusBadRequest,
		},
		"validation: name too short": {
			reqBody:   `{"name": "A"}`,
			resStatus: http.StatusBadRequest,
		},
		"already exists": {
			reqBody:   `{"name": "UserA-Name"}`,
			dbUsers:   []User{tfUserA},
			resStatus: http.StatusBadRequest,
		},
		"JSON: malformed": {
			reqBody:   `NotA-JSON`,
			resStatus: http.StatusBadRequest,
		},
		"JSON: type not matching": {
			reqBody:   `[{"A":1},{"B":"123"}]`,
			resStatus: http.StatusBadRequest,
		},
		"UserSave error": {
			reqBody:     `{"name": "ABCDE"}`,
			usErr:       errors.New("some kind of DB error"),
			usCalledExp: true,
			resStatus:   http.StatusInternalServerError,
		},
	}

	for sym, tc := range tests {
		st := NewTmMemoryStorageMock()
		st.outUserSaveErr = tc.usErr

		h := NewHTTPDefaultHandler(st)
		ts = httptest.NewServer(h)

		// GIVEN: expected users are in DB
		for _, u := range tc.dbUsers {
			uC := u
			ar.NoError(t, st.UserSave(&uC), "case: %s", sym)
		}
		st.inUserSaveCalled = false

		// WHEN: user create is called
		res, err := http.Post(fmt.Sprintf("%s/v1/users", ts.URL), "application/json", strings.NewReader(tc.reqBody))

		// THEN: validate response
		ar.NoError(t, err, "[%s] unexpected error from HTTP client", sym)
		a.EqualValues(t, 0, res.ContentLength, "[%s] non empty response body", sym)
		a.Equal(t, tc.resStatus, res.StatusCode, "[%s] mismatch on response code", sym)

		// AND: validate storage access
		a.Equal(t, tc.usCalledExp, st.inUserSaveCalled, "[%s] userSave function call status mismatch", sym)
	}
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

	matches := rPathMsgRead.FindStringSubmatch(res.Header.Get("Location"))
	// matches also have the source string on index 0
	ar.Len(t, matches, 2, "response: location header does not point to message read action")
	msgIDFromHeader := matches[1]

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
	a.Equal(t, msgIDFromHeader, msgGot.ID, "message ID mismatch")
}

func Test_HTTPHandler_Message_Create_Failure(t *testing.T) {
	var ts *httptest.Server
	defer func() {
		if ts != nil {
			ts.Close()
		}
	}()

	tests := map[string]struct {
		reqBody     string
		dbUsers     []User
		ufErr       error // uf = UserFind
		ufCalledExp bool
		msErr       error // ms = MsgSave
		msCalledExp bool
		resStatus   int
	}{
		"JSON: malformed": {
			reqBody:   `NotA-JSON`,
			resStatus: http.StatusBadRequest,
		},
		"JSON: type not matching": {
			reqBody:   `[{"A":1},{"B":"123"}]`,
			resStatus: http.StatusBadRequest,
		},
		"validation: empty object/name empty": {
			reqBody:   `{}`,
			resStatus: http.StatusBadRequest,
		},
		"unknown author": {
			reqBody:     `{"author": "UserUnknown-Name","body":"qweasd","tag":"tagA"}`,
			dbUsers:     []User{tfUserA},
			ufCalledExp: true,
			resStatus:   http.StatusBadRequest,
		},
		"MsgSave error": {
			reqBody:     tfTrInMsgAA_JSON,
			dbUsers:     []User{tfUserA},
			ufCalledExp: true,
			msErr:       errors.New("some kind of DB error"),
			msCalledExp: true,
			resStatus:   http.StatusInternalServerError,
		},
		"UserFind error": {
			reqBody:     tfTrInMsgAA_JSON,
			ufCalledExp: true,
			ufErr:       errors.New("some kind of DB error"),
			msCalledExp: false,
			resStatus:   http.StatusInternalServerError,
		},
	}

	for sym, tc := range tests {
		st := NewTmMemoryStorageMock()
		st.outUserFindErr = tc.ufErr
		st.outMsgSaveErr = tc.msErr

		h := NewHTTPDefaultHandler(st)
		ts = httptest.NewServer(h)

		// GIVEN: expected users are in DB
		for _, u := range tc.dbUsers {
			uC := u
			ar.NoError(t, st.UserSave(&uC), "case: %s", sym)
		}
		st.inUserSaveCalled = false

		// WHEN: message create is called
		res, err := http.Post(fmt.Sprintf("%s/v1/messages", ts.URL), "application/json", strings.NewReader(tc.reqBody))

		// THEN: validate response
		ar.NoError(t, err, "[%s] unexpected error from HTTP client", sym)
		a.EqualValues(t, 0, res.ContentLength, "[%s] non empty response body", sym)
		a.Equal(t, tc.resStatus, res.StatusCode, "[%s] mismatch on response code", sym)

		// AND: validate storage access
		a.Equal(t, tc.ufCalledExp, st.inUserFindCalled, "[%s] UserFind function call status mismatch", sym)
		a.Equal(t, tc.msCalledExp, st.inMsgSaveCalled, "[%s] MsgSave function call status mismatch", sym)
	}
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
		exp     MessagesCollectionOut
	}{
		"single": {
			[]User{tfUserA},
			[]Message{tfMsgAA},
			tfTrInMsgAA.Tag,
			MessagesCollectionOut{tfTrOutMsgAA},
		},
		"two": {
			[]User{tfUserA},
			[]Message{tfMsgAA, tfMsgAB},
			tfTrInMsgAA.Tag,
			MessagesCollectionOut{tfTrOutMsgAA, tfTrOutMsgAB},
		},
		"three, cross user, partial": {
			[]User{tfUserA, tfUserB},
			[]Message{tfMsgAA, tfMsgAB, tfMsgBA, tfMsgBB},
			tfTrInMsgAA.Tag,
			MessagesCollectionOut{tfTrOutMsgAA, tfTrOutMsgAB, tfTrOutMsgBA},
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
		a.Equal(t, "application/json", res.Header.Get("Content-Type"), "mismatch on response content encoding")

		var resBodyGot MessagesCollectionOut
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
}

func Test_HTTPHandler_Message_Find_Failure(t *testing.T) {
	var ts *httptest.Server
	defer func() {
		if ts != nil {
			ts.Close()
		}
	}()

	tests := map[string]struct {
		dbUsers     []User
		dbMsg       []Message
		tag         Tag
		mfCalledExp bool // mf = MsgFind
		mfErr       error
		mlCalledExp bool // ml = MsgLoad
		mlErr       error
		ulCalledExp bool // ul = UserLoad
		ulErr       error
		resStatus   int
	}{
		"MsgFind error": {
			dbUsers:     []User{tfUserA},
			dbMsg:       []Message{tfMsgAA},
			tag:         tfTrInMsgAA.Tag,
			mfCalledExp: true,
			mfErr:       errors.New("some kind of DB error"),
			resStatus:   http.StatusInternalServerError,
		},
		"MsgLoad error": {
			dbUsers:     []User{tfUserA},
			dbMsg:       []Message{tfMsgAA},
			tag:         tfTrInMsgAA.Tag,
			mfCalledExp: true,
			mlCalledExp: true,
			mlErr:       errors.New("some kind of DB error"),
			resStatus:   http.StatusInternalServerError,
		},
		"UserLoad error": {
			dbUsers:     []User{tfUserA},
			dbMsg:       []Message{tfMsgAA},
			tag:         tfTrInMsgAA.Tag,
			mfCalledExp: true,
			mlCalledExp: true,
			ulCalledExp: true,
			ulErr:       errors.New("some kind of DB error"),
			resStatus:   http.StatusInternalServerError,
		},
	}

	for sym, tc := range tests {
		st := NewTmMemoryStorageMock()
		st.outMsgFindErr = tc.mfErr
		st.outMsgLoadErr = tc.mlErr
		st.outUserLoadErr = tc.ulErr

		h := NewHTTPDefaultHandler(st)
		ts = httptest.NewServer(h)

		// GIVEN: user and message are in DB
		for _, u := range tc.dbUsers {
			uC := u
			ar.NoError(t, st.UserSave(&uC), "case: %s", sym)
		}
		st.inUserSaveCalled = false
		for _, m := range tc.dbMsg {
			mC := m
			ar.NoError(t, st.MsgSave(&mC), "case: %s", sym)
		}
		st.inMsgSaveCalled = false

		// WHEN: message create is called
		res, err := http.Get(fmt.Sprintf("%s/v1/messages?tag=%s", ts.URL, tc.tag))

		// THEN: validate response
		ar.NoError(t, err, "[%s] unexpected error from HTTP client", sym)
		a.EqualValues(t, 0, res.ContentLength, "[%s] non empty response body", sym)
		a.Equal(t, tc.resStatus, res.StatusCode, "[%s] mismatch on response code", sym)

		// AND: validate storage access
		a.Equal(t, tc.mfCalledExp, st.inMsgFindCalled, "[%s] MsgFind function call status mismatch", sym)
		a.Equal(t, tc.mlCalledExp, st.inMsgLoadCalled, "[%s] MsgLoad function call status mismatch", sym)
		a.Equal(t, tc.ulCalledExp, st.inUserLoadCalled, "[%s] UserLoad function call status mismatch", sym)
	}
}

func Test_HTTPHandler_Message_Read_Success_Found(t *testing.T) {
	st := NewTmMemoryStorageMock()
	h := NewHTTPDefaultHandler(st)
	ts := httptest.NewServer(h)
	defer ts.Close()

	// GIVEN: matching message is in DB
	user := tfUserA
	msg := tfMsgAA
	ar.NoError(t, st.UserSave(&user))
	ar.NoError(t, st.MsgSave(&msg))

	res, err := http.Get(fmt.Sprintf("%s/v1/messages/%s", ts.URL, msg.ID))

	// THEN: validate response
	ar.NoError(t, err, "unexpected error from HTTP client")

	a.Equal(t, http.StatusOK, res.StatusCode, "mismatch on response code")
	a.Equal(t, "application/json", res.Header.Get("Content-Type"), "mismatch on response content encoding")

	var resBodyGot MessageOut
	err = json.NewDecoder(res.Body).Decode(&resBodyGot)
	res.Body.Close()
	ar.NoError(t, err, "unexpected error on response body read")

	a.Equal(t, tfTrOutMsgAA, resBodyGot)
}

func Test_HTTPHandler_Message_Read_Success_NotFound(t *testing.T) {
	st := NewMemoryStorage()
	h := NewHTTPDefaultHandler(st)
	ts := httptest.NewServer(h)
	defer ts.Close()

	// GIVEN: message NOT in DB
	res, err := http.Get(fmt.Sprintf("%s/v1/messages/non-existing-123", ts.URL))

	// THEN: validate response
	ar.NoError(t, err, "unexpected error from HTTP client")
	a.Equal(t, http.StatusNotFound, res.StatusCode, "mismatch on response code")
	a.EqualValues(t, 0, res.ContentLength, "non empty response body")
}

func Test_HTTPHandler_Message_Read_Failure(t *testing.T) {
	var ts *httptest.Server
	defer func() {
		if ts != nil {
			ts.Close()
		}
	}()

	tests := map[string]struct {
		dbUsers     []User
		dbMsg       []Message
		msgID       string
		mlCalledExp bool // ml = MsgLoad
		mlErr       error
		ulCalledExp bool // ul = UserLoad
		ulErr       error
		resStatus   int
	}{
		"MsgLoad error": {
			dbUsers:     []User{tfUserA},
			dbMsg:       []Message{tfMsgAA},
			msgID:       tfMsgAA.ID,
			mlCalledExp: true,
			mlErr:       errors.New("some kind of DB error"),
			resStatus:   http.StatusInternalServerError,
		},
		"UserLoad error: user not found": {
			dbUsers:     []User{tfUserA},
			dbMsg:       []Message{tfMsgAA},
			msgID:       tfMsgAA.ID,
			mlCalledExp: true,
			ulCalledExp: true,
			ulErr:       ErrElementNotFound,
			resStatus:   http.StatusInternalServerError,
		},
		"UserLoad error: db error": {
			dbUsers:     []User{tfUserA},
			dbMsg:       []Message{tfMsgAA},
			msgID:       tfMsgAA.ID,
			mlCalledExp: true,
			ulCalledExp: true,
			ulErr:       errors.New("some kind of DB error"),
			resStatus:   http.StatusInternalServerError,
		},
	}

	for sym, tc := range tests {
		st := NewTmMemoryStorageMock()
		st.outMsgLoadErr = tc.mlErr
		st.outUserLoadErr = tc.ulErr

		h := NewHTTPDefaultHandler(st)
		ts = httptest.NewServer(h)

		// GIVEN: user and message are in DB
		for _, u := range tc.dbUsers {
			uC := u
			ar.NoError(t, st.UserSave(&uC), "case: %s", sym)
		}
		st.inUserSaveCalled = false
		for _, m := range tc.dbMsg {
			mC := m
			ar.NoError(t, st.MsgSave(&mC), "case: %s", sym)
		}
		st.inMsgSaveCalled = false

		// WHEN: message create is called
		res, err := http.Get(fmt.Sprintf("%s/v1/messages/%s", ts.URL, tc.msgID))

		// THEN: validate response
		ar.NoError(t, err, "[%s] unexpected error from HTTP client", sym)
		a.EqualValues(t, 0, res.ContentLength, "[%s] non empty response body", sym)
		a.Equal(t, tc.resStatus, res.StatusCode, "[%s] mismatch on response code", sym)

		// AND: validate storage access
		a.Equal(t, tc.mlCalledExp, st.inMsgLoadCalled, "[%s] MsgLoad function call status mismatch", sym)
		a.Equal(t, tc.ulCalledExp, st.inUserLoadCalled, "[%s] UserLoad function call status mismatch", sym)
	}
}

func Test_HTTPHandler_Message_GET_unknownPath(t *testing.T) {
	st := NewMemoryStorage()
	h := NewHTTPDefaultHandler(st)
	ts := httptest.NewServer(h)
	defer ts.Close()

	// GIVEN: message NOT in DB
	res, err := http.Get(fmt.Sprintf("%s/v1/messages/some-kind-of/strange-123-path", ts.URL))

	// THEN: validate response
	ar.NoError(t, err, "unexpected error from HTTP client")
	a.Equal(t, http.StatusNotFound, res.StatusCode, "mismatch on response code")
	a.EqualValues(t, 0, res.ContentLength, "non empty response body")
}

// TODO: validate file content
func Test_HTTPHandler_Swagger(t *testing.T) {
	st := NewMemoryStorage()
	h := NewHTTPDefaultHandler(st)
	ts := httptest.NewServer(h)
	defer ts.Close()

	res, err := http.Get(fmt.Sprintf("%s/v1/swagger.json", ts.URL))

	ar.NoError(t, err, "unexpected error from HTTP client")
	a.Equal(t, http.StatusOK, res.StatusCode, "mismatch on response code")
	a.NotZero(t, res.ContentLength, "empty response body")
}

func Test_HTTPHandler_Options(t *testing.T) {
	var ts *httptest.Server
	defer func() {
		if ts != nil {
			ts.Close()
		}
	}()

	tests := []struct {
		path string
	}{
		{"/v1/users"},
	}

	for _, tc := range tests {
		st := NewTmMemoryStorageMock()

		h := NewHTTPDefaultHandler(st)
		ts = httptest.NewServer(h)

		req, err := http.NewRequest(http.MethodOptions, fmt.Sprintf("%s%s", ts.URL, tc.path), nil)
		ar.NoError(t, err, "[%s] unexpected error from request creation", tc.path)

		res, err := http.DefaultClient.Do(req)
		ar.NoError(t, err, "[%s] unexpected error from call", tc.path)

		a.EqualValues(t, 0, res.ContentLength, "[%s] non empty response body", tc.path)
		a.Equal(t, http.StatusOK, res.StatusCode, "[%s] mismatch on response code", tc.path)
	}
}
