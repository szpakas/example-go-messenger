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
	AuthorID: tfUserA.ID,
	Tag:      tfTagA,
}

var tfMsgAB = Message{
	ID:       "UserA_MessageB-ID",
	Body:     "UserA_MessageB-Body",
	AuthorID: tfUserA.ID,
	Tag:      tfTagA,
}

var tfMsgAXA_NoID = Message{
	Body:     "UserA_MessageXA-Body",
	AuthorID: tfUserA.ID,
	Tag:      tfTagA,
}

var tfMsgAXB_NoBody = Message{
	ID:       "MessageAXB-Body",
	AuthorID: tfUserA.ID,
	Tag:      tfTagA,
}

var tfMsgXXA_NoAuthorID = Message{
	ID:   "MessageXXA-Body",
	Body: "MessageXXA-Body",
	Tag:  tfTagA,
}

var tfMsgAXC_NoTag = Message{
	ID:       "MessageAXC-Body",
	Body:     "UserA_MessageXC-Body",
	AuthorID: tfUserA.ID,
}

var tfMsgAXD_TagTooShort = Message{
	ID:       "MessageAXD-Body",
	Body:     "UserA_MessageXD-Body",
	AuthorID: tfUserA.ID,
	Tag:      tfTagXA_TooShort,
}

var tfMsgBA = Message{
	ID:       "UserB_MessageA-ID",
	Body:     "UserB_MessageA-Body",
	AuthorID: tfUserB.ID,
	Tag:      tfTagA,
}

var tfMsgBB = Message{
	ID:       "UserB_MessageB-ID",
	Body:     "UserB_MessageB-Body",
	AuthorID: tfUserB.ID,
	Tag:      tfTagB,
}

// -- section: Tag
var tfTagA = Tag("tagA")
var tfTagB = Tag("tagB")
var tfTagC = Tag("tagC")

var tfTagXA_TooShort = Tag("s")
