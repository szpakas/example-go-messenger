package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	"github.com/satori/go.uuid"
)

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

	// swagger:route POST /v1/users users UserCreate
	//
	// Create user.
	//
	//     Responses:
	//       201: UserCreatedResponse
	//       400: BadRequestError
	//       500: InternalServerError
	mux.Handle("/v1/users", &usersHandler{Storer: st})

	mux.Handle("/v1/messages", &messagesHandler{Storer: st})
	// duplication needed to handle base path without redirection
	mux.Handle("/v1/messages/", &messagesHandler{Storer: st})

	mux.Handle("/v1/swagger.json", &swaggerHandler{})

	return mux
}

// usersHandler is HTTP handler for users related actions
type usersHandler struct {
	Storer UserStorer
}

func (h *usersHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	var trIn UserIn
	err := json.NewDecoder(r.Body).Decode(&trIn)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if trIn.Validate() != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if _, err := h.Storer.UserFindByName(trIn.Name); err != ErrElementNotFound {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user := User{
		ID:   uuid.NewV1().String(),
		Name: trIn.Name,
	}

	err = h.Storer.UserSave(&user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// messagesHandler is HTTP handler for messages related actions
type messagesHandler struct {
	Storer Storer
}

func (h *messagesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch true {
	case r.Method == http.MethodPost:
		// swagger:route POST /v1/messages messages MessageCreate
		//
		// Create message.
		//
		//     Responses:
		//       201: MessageCreatedResponse
		//       400: BadRequestError
		//       500: InternalServerError
		h.handleCreate(w, r)
		return
	case r.Method == http.MethodGet && (r.URL.Path == "/v1/messages" || r.URL.Path == "/v1/messages/"):
		// swagger:route GET /v1/messages messages MessagesFind
		//
		// Get collection of messages matching requested tag.
		//
		//     Responses:
		//       200: MessagesCollectionResponse
		//       404: NotFoundError
		//       500: InternalServerError
		h.handleFind(w, r)
		return
	case r.Method == http.MethodGet:
		// swagger:route GET /v1/messages/{id} messages MessageRead
		//
		// Get details of single message by its ID.
		//
		//     Responses:
		//       200: MessageReadResponse
		//       404: NotFoundError
		//       500: InternalServerError
		h.handleRead(w, r)
		return
	}
}

func (h *messagesHandler) handleCreate(w http.ResponseWriter, r *http.Request) {
	var trIn MessageIn
	err := json.NewDecoder(r.Body).Decode(&trIn)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if trIn.Validate() != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	author, err := h.Storer.UserFindByName(trIn.Author)
	switch err {
	case nil:
	case ErrElementNotFound:
		w.WriteHeader(http.StatusBadRequest)
		return
	default:
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	msg := Message{
		ID:       uuid.NewV1().String(),
		Body:     trIn.Body,
		Tag:      trIn.Tag,
		AuthorID: author.ID,
	}

	err = h.Storer.MsgSave(&msg)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Location", "/v1/messages/"+msg.ID)
	w.WriteHeader(http.StatusCreated)
}

func msgToTransport(msg *Message, author *User) MessageOut {
	return MessageOut{
		ID:     msg.ID,
		Author: author.Name,
		Body:   msg.Body,
		Tag:    msg.Tag,
	}
}

func (h *messagesHandler) handleFind(w http.ResponseWriter, r *http.Request) {
	tag := r.URL.Query().Get("tag")
	msgsIDs, err := h.Storer.MsgsIDsFindByTag(Tag(tag))

	switch err {
	case nil:
	case ErrElementNotFound:
		w.WriteHeader(http.StatusNotFound)
		return
	default:
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var trOut MessagesCollectionOut
	for _, mID := range msgsIDs {
		msg, err := h.Storer.MsgLoad(mID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		author, err := h.Storer.UserLoad(msg.AuthorID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		trOut = append(trOut, msgToTransport(msg, author))
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(trOut)
}

// allowed chars in ID: 0-9a-zA-Z-_ (space is NOT allowed)
var rPathMsgRead = regexp.MustCompile(`^/v1/messages/([\da-zA-Z\-_]+)/?$`)

func (h *messagesHandler) handleRead(w http.ResponseWriter, r *http.Request) {
	matches := rPathMsgRead.FindStringSubmatch(r.URL.Path)

	// matches also have the source string on index 0
	if len(matches) != 2 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// msgID is on index 1
	msg, err := h.Storer.MsgLoad(matches[1])
	switch err {
	case nil:
	case ErrElementNotFound:
		w.WriteHeader(http.StatusNotFound)
		return
	default:
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	author, err := h.Storer.UserLoad(msg.AuthorID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	trOut := msgToTransport(msg, author)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(trOut)
}

// swaggerHandler is HTTP handler for swagger definition file
type swaggerHandler struct{}

func (h *swaggerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./swagger.json")
}
