// Package classification Messenger API.
//
// the purpose of this application is to provide an test bed platform
// for toying with GO
//
// Terms Of Service:
//
// there are no TOS at this moment, use at your own risk we take no responsibility
//
//     Schemes: http
//     Host: localhost:8080
//     BasePath: /
//     Version: 0.0.1
//     License: Apache 2.0
//     Contact: Adam Szpakowski <szpakas@gmail.com>
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
// swagger:meta
package main

// A UserBodyParams model.
//
// This is used for operations that want an User as body of the request
//
// swagger:parameters UserCreate
type UserBodyParams struct {
	// User to create
	//
	// in: body
	// required: true
	User *UserIn `json:"user"`
}

// UserCreatedResponse represents response to creation of the user.
//
// swagger:response UserCreatedResponse
type UserCreatedResponse struct{}

// A MessageBodyParams model.
//
// This is used for operations that want an Message as body of the request
//
// swagger:parameters MessageCreate
type MessageBodyParams struct {
	// The message to submit
	//
	// in: body
	// required: true
	Message *MessageIn `json:"message"`
}

// A MessageQueryFlags contains the query flags for things that list tags
//
// swagger:parameters MessagesFind
type MessageQueryFlags struct {
	// Tag attached to the message
	//
	// in: query
	// required: true
	Tag string `json:"tag"`
}

// MessageCreatedResponse represents response to creation of the message.
//
// swagger:response MessageCreatedResponse
type MessageCreatedResponse struct {
	// Location is relative URL to newly created user.
	Location string
}

// MessageResponse represents transport level model for single message returned from system to user.
//
// swagger:response MessageReadResponse
type MessageReadResponse struct {
	// in: body
	Body *MessageOut
}

// MessagesCollectionResponse represents transport level model for collection of messages returned from system to user.
//
// swagger:response MessagesCollectionResponse
type MessagesCollectionResponse struct {
	// in: body
	Body []*MessageOut
}

// A BadRequestError is an error that is generated when user submitted request which is incorrect.
// One of the cases is some kind of validation error.
// Repeating the request will most probably not change the outcome.
//
// swagger:response BadRequestError
type BadRequestError struct{}

// A NotFoundError is an error that is generated when requested element could not be found.
// It's also used when collection is requested for specific parameters combination and returns empty.
//
// swagger:response NotFoundError
type NotFoundError struct{}

// A InternalServerError is an error that is generated when server could not produce response.
// Repeating the request will most probably not change the outcome.
//
// swagger:response InternalServerError
type InternalServerError struct{}
