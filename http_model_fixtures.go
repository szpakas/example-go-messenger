package main

// -- section: User
var tfTrInUserA = UserIn{
	Name: "UserA-Name",
}
var tfTrInUserA_JSON = `{"name":"UserA-Name"}`

var tfTrInUserB = UserIn{
	Name: "UserB-Name",
}

// -- section: Message
var tfTrInMsgAA = MessageIn{
	Body:   "UserA_MessageA-Body",
	Author: "UserA-Name",
	Tag:    Tag("tagA"),
}
var tfTrInMsgAA_JSON = `{"body":"UserA_MessageA-Body","author":"UserA-Name","tag":"tagA"}`

var tfTrInMsgAXB_NoBody = MessageIn{
	Author: tfUserA.Name,
	Tag:    tfTagA,
}

var tfTrInMsgXXA_NoAuthor = MessageIn{
	Body: "MessageXXA-Body",
	Tag:  tfTagA,
}

var tfTrInMsgAXC_NoTag = MessageIn{
	Body:   "UserA_MessageXC-Body",
	Author: tfUserA.Name,
}

var tfTrInMsgAXD_TagTooShort = MessageIn{
	Body:   "UserA_MessageXD-Body",
	Author: tfUserA.Name,
	Tag:    tfTagXA_TooShort,
}

var tfTrOutMsgAA = MessageOut{
	ID:     "UserA_MessageA-ID",
	Body:   "UserA_MessageA-Body",
	Author: "UserA-Name",
	Tag:    Tag("tagA"),
}

var tfTrOutMsgAA_JSON = `{"id":"UserA_MessageA-ID","body":"UserA_MessageA-Body","author":"UserA-Name","tag":"tagA"}`

var tfTrOutMsgAB = MessageOut{
	ID:     "UserA_MessageB-ID",
	Body:   "UserA_MessageB-Body",
	Author: "UserA-Name",
	Tag:    Tag("tagA"),
}

var tfTrOutMsgAB_JSON = `{"id":"UserA_MessageB-ID","body":"UserA_MessageB-Body","author":"UserA-Name","tag":"tagA"}`

var tfTrOutMsgBA = MessageOut{
	ID:     "UserB_MessageA-ID",
	Body:   "UserB_MessageA-Body",
	Author: "UserB-Name",
	Tag:    Tag("tagA"),
}
