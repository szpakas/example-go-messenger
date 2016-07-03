package main

var (
	userNameLengthMin = 2
)

// UserIn represents transport level model for single user submitted into the HTTP handler.
type UserIn struct {
	// Name represents the user to the outside world.
	//
	// required: true
	// min length: 2
	Name string `json:"name"`
}

// Validate validates the User and returns error on failure.
func (u UserIn) Validate() error {
	if u.Name == "" {
		return NewValidationError("Name missing")
	}
	if len(u.Name) < userNameLengthMin {
		return NewValidationError("Name too short")
	}
	return nil
}

// MessageIn represents transport level model for single message sent by user to the system.
type MessageIn struct {
	// Body represents the actual message
	//
	// required: true
	Body string `json:"body"`

	// Author is an Name of the user who authored message
	//
	// required: true
	Author string `json:"author"`

	// Tag is a tag attached to a message
	//
	// required: true
	// min length: 2
	Tag Tag `json:"tag"`
}

// Validate validates the Message and returns error on failure.
func (m MessageIn) Validate() error {
	if m.Body == "" {
		return NewValidationError("missing Body")
	}
	if m.Author == "" {
		return NewValidationError("missing Author")
	}
	if err := m.Tag.Validate(); err != nil {
		return NewValidationError("invalid Tag", err)
	}
	return nil
}

// A MessageID parameter model.
//
// This is used for operations that want the ID of an message in the path
//
// swagger:parameters MessageRead
type MessageID struct {
	// ID represents the unique identifier for the message
	//
	// in: path
	// required: true
	ID string `json:"id"` // json tag is used to modify the swagger naming
}

type MessageOut struct {
	// ID represents the unique identifier for the message
	//
	// required: true
	ID string `json:"id"`

	// Body represents the actual message
	//
	// required: true
	Body string `json:"body"`

	// Author is an Name of the user who authored message
	//
	// required: true
	Author string `json:"author"`

	// Tag is a tag attached to a message
	//
	// required: true
	Tag Tag `json:"tag"`
}

type MessagesCollectionOut []MessageOut
