package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/satori/go.uuid"

	"github.com/davecgh/go-spew/spew"
)

var _ = spew.Config

// UserStorer is storage interface for User related operations
type UserStorer interface {
	UserSave(u *User) error
	UserLoad(id string) (*User, error)
	UserFindByName(name string) (*User, error)
}

// MsgStorer is storage interface for Message related operations
type MsgStorer interface {
	MsgSave(m *Message) error
	MsgLoad(id string) (*Message, error)
	MsgsIDsFindByTag(tag Tag) ([]string, error)
}

// Storer is an storage interface for users, messages and tags
type Storer interface {
	UserStorer
	MsgStorer
}

// NewHTTPServer creates new HTTP server for package submission.
func NewHTTPServer(host string, port int, h http.Handler) *http.Server {
	return &http.Server{
		Addr:    fmt.Sprintf("%s:%d", host, port),
		Handler: h,
	}
}

// NewHTTPDefaultHandler is a default handler factory.
// It takes care of routing.
// TODO: test me
func NewHTTPDefaultHandler(st Storer) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/v1/users", &usersHandler{Storer: st})
	mux.Handle("/v1/messages", &messagesHandler{Storer: st})

	return mux
}

// usersHandler is HTTP handler for users related actions
type usersHandler struct {
	Storer UserStorer
}

func (h *usersHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var trIn TrInUser
	_ = json.NewDecoder(r.Body).Decode(&trIn)

	user := User{
		ID:   uuid.NewV1().String(),
		Name: trIn.Name,
	}

	_ = h.Storer.UserSave(&user)

	w.WriteHeader(http.StatusCreated)
}

// messagesHandler is HTTP handler for users related actions
type messagesHandler struct {
	Storer Storer
}

func (h *messagesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.handleCreate(w, r)
		return
	case http.MethodGet:
		h.handleFind(w, r)
		return
	}
}

func (h *messagesHandler) handleCreate(w http.ResponseWriter, r *http.Request) {
	var trIn TrInMessage
	_ = json.NewDecoder(r.Body).Decode(&trIn)

	author, _ := h.Storer.UserFindByName(trIn.Author)

	msg := Message{
		ID:       uuid.NewV1().String(),
		Body:     trIn.Body,
		Tag:      trIn.Tag,
		AuthorID: author.ID,
	}

	_ = h.Storer.MsgSave(&msg)

	w.WriteHeader(http.StatusCreated)
}

func (h *messagesHandler) handleFind(w http.ResponseWriter, r *http.Request) {
	tag := r.URL.Query().Get("tag")
	msgsIDs, err := h.Storer.MsgsIDsFindByTag(Tag(tag))

	switch err {
	case nil:
	case ErrElementNotFound:
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var trOut TrOutMessagesCollection
	for _, mID := range msgsIDs {
		msg, _ := h.Storer.MsgLoad(mID)
		user, _ := h.Storer.UserLoad(msg.AuthorID)

		trOut = append(trOut, TrOutMessage{
			ID:     msg.ID,
			Author: user.Name,
			Body:   msg.Body,
			Tag:    msg.Tag,
		})
	}

	w.Header().Set("Content-Encoding", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(trOut)
}
