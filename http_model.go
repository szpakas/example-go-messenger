package main

var (
	userNameLengthMin = 2
)

// TrInUser represents transport level model for single user submitted into the HTTP handler.
type TrInUser struct {
	// Name represents the user to the outside world.
	Name string `json:"name"`
}

// Validate validates the User and returns error on failure.
func (m TrInUser) Validate() error {
	if m.Name == "" {
		return NewValidationError("Name missing")
	}
	if len(m.Name) < userNameLengthMin {
		return NewValidationError("Name too short")
	}
	return nil
}

// TrInMessage represents transport level model for single message sent by user to the system.
type TrInMessage struct {
	// Body represents the actual message.
	Body string `json:"body"`

	// Author is an Name of the user who authored message.
	Author string `json:"author"`

	// Tag is a tag attached to a message
	Tag Tag `json:"tag"`
}

// Validate validates the Message and returns error on failure.
func (m TrInMessage) Validate() error {
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

// TrOutMessage represents transport level model for single message returned from system to user.
type TrOutMessage struct {
	// Body represents the actual message.
	ID string `json:"id"`

	// Body represents the actual message.
	Body string `json:"body"`

	// Author is an Name of the user who authored message.
	Author string `json:"author"`

	// Tag is a tag attached to a message
	Tag Tag `json:"tag"`
}

// TrOutMessagesCollection represents collection of messages returned from server
type TrOutMessagesCollection []TrOutMessage
