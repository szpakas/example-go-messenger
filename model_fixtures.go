package main

// -- section: User
var tfUserA = User{
	ID:   "UserA-ID",
	Name: "UserA-Name",
}

var tfUserB = User{
	ID:   "UserB-ID",
	Name: "UserB-Name",
}

var tfUserXA_NoID = User{
	Name: "UserXA-Name",
}

// -- section: Message
var tfMsgAA = Message{
	ID:       "UserA_MessageA-ID",
	Body:     "UserA_MessageA-Body",
	AuthorID: "UserA-ID",
	Tag:      Tag("tagA"),
}

var tfMsgAB = Message{
	ID:       "UserA_MessageB-ID",
	Body:     "UserA_MessageB-Body",
	AuthorID: "UserA-ID",
	Tag:      Tag("tagA"),
}

var tfMsgBA = Message{
	ID:       "UserB_MessageA-ID",
	Body:     "UserB_MessageA-Body",
	AuthorID: "UserB-ID",
	Tag:      Tag("tagA"),
}

var tfMsgBB = Message{
	ID:       "UserB_MessageB-ID",
	Body:     "UserB_MessageB-Body",
	AuthorID: "UserB-ID",
	Tag:      Tag("tagB"),
}

var tfMsgAXA_NoID = Message{
	Body:     "UserA_MessageXA-Body",
	AuthorID: "UserA-ID",
	Tag:      Tag("tagB"),
}

// -- section: Tag
var tfTagA = Tag("tagA")
var tfTagB = Tag("tagB")
var tfTagC = Tag("tagC")

var tfTagXA_TooShort = Tag("s")
