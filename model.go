package main

var (
	tagLengthMin = 2
	tagLengthMax = 128
)

// User represents model for single user using the system.
type User struct {
	// ID is a unique, immutable identifier for the user.
	// Used in internal relations, logging.
	ID string

	// Name represents the user to the outside world.
	// It may be changed and shall never be used for anything else then human interaction.
	Name string
}

// Message represents model for single message sent by user to the system.
type Message struct {
	// ID is a unique, immutable identifier for the message.
	// Used in internal relations, logging.
	ID string

	// Body represents the actual message.
	Body string

	// AuthorID is an ID of the user who authored message.
	AuthorID string

	// Tag is a tag attached to a message
	Tag Tag
}

// Tag represents model for a single Tag attached to a message.
type Tag string

// Validate validates the tag and returns error on failure.
func (t Tag) Validate() error {
	if t == "" {
		return NewValidationError("empty value")
	}
	if len(t) < tagLengthMin {
		return NewValidationError("too short")
	}
	if len(t) > tagLengthMax {
		return NewValidationError("too long")
	}
	return nil
}
