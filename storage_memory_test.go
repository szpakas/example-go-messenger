package main

import (
	"testing"

	a "github.com/stretchr/testify/assert"
	ar "github.com/stretchr/testify/require"
)

func Test_MemoryStorage_Factory(t *testing.T) {
	s, closer := tsMemoryStorageSetup()
	defer closer()

	ar.NotNil(t, s, "empty element returned")
	ar.IsType(t, &memoryStorage{}, s)

	a.NotZero(t, s.users, "users map is not initialised")
	a.Len(t, s.users, 0, "users map should be empty on init")

	a.NotZero(t, s.messages, "messages map is not initialised")
	a.Len(t, s.messages, 0, "messages map should be empty on init")

	a.NotZero(t, s.tags, "tags map is not initialised")
	a.Len(t, s.tags, 0, "tags map should be empty on init")
}

// -- section: User
func Test_MemoryStorage_UserSave_Success(t *testing.T) {
	s, closer := tsMemoryStorageSetup()
	defer closer()

	elExp := tfUserA
	ar.NoError(t, s.UserSave(&elExp))
	s.usersMu.RLock()
	defer s.usersMu.RUnlock()
	ar.Contains(t, s.users, elExp.ID, "User with requested ID is not in storage")
	a.Equal(t, s.users[elExp.ID], &elExp, "User from storage does not match")
}

func Test_MemoryStorage_UserSave_Failure_NoID(t *testing.T) {
	s, closer := tsMemoryStorageSetup()
	defer closer()

	ar.EqualError(t, s.UserSave(&tfUserXA_NoID), ErrElementIDNotSet.Error())
	s.usersMu.RLock()
	defer s.usersMu.RUnlock()
	ar.Len(t, s.users, 0, "unexpected element stored")
}

func Test_MemoryStorage_UserLoad_Exists(t *testing.T) {
	s, closer := tsMemoryStorageSetup()
	defer closer()

	elExp := tfUserA

	// GIVEN: expected user is in storage
	ar.NoError(t, s.UserSave(&elExp))

	elGot, err := s.UserLoad(elExp.ID)
	ar.NoError(t, err)
	a.Equal(t, &elExp, elGot, "User from storage does not match")
}

func Test_MemoryStorage_UserLoad_NotFound(t *testing.T) {
	s, closer := tsMemoryStorageSetup()
	defer closer()

	// GIVEN: expected user is NOT in storage

	_, err := s.UserLoad(tfUserA.ID)
	ar.Equal(t, ErrElementNotFound, err)
}

func Test_MemoryStorage_UserFindByName_Exists(t *testing.T) {
	s, closer := tsMemoryStorageSetup()
	defer closer()

	elExp := tfUserA

	// GIVEN: expected user is in storage
	ar.NoError(t, s.UserSave(&elExp))

	elGot, err := s.UserFindByName(elExp.Name)
	ar.NoError(t, err)
	a.Equal(t, &elExp, elGot, "User from storage does not match")
}

func Test_MemoryStorage_UserFindByName_NotFound(t *testing.T) {
	s, closer := tsMemoryStorageSetup()
	defer closer()

	// GIVEN: expected user is NOT in storage

	_, err := s.UserFindByName(tfUserA.Name)
	ar.Equal(t, ErrElementNotFound, err)
}

// -- section: Message
func Test_MemoryStorage_MessageSave_Success(t *testing.T) {
	s, closer := tsMemoryStorageSetup()
	defer closer()

	// GIVEN: expected user is in storage
	userExp := tfUserA
	ar.NoError(t, s.UserSave(&userExp))

	msgExp := tfMsgAA
	ar.NoError(t, s.MsgSave(&msgExp))

	// THEN: messages storage has message
	s.messagesMu.RLock()
	defer s.messagesMu.RUnlock()
	ar.Contains(t, s.messages, msgExp.ID, "Message with requested ID is not in storage")
	a.Equal(t, s.messages[msgExp.ID], &msgExp, "Message from storage does not match")

	// AND: tag is mapped
	s.tagsMu.RLock()
	defer s.tagsMu.RUnlock()
	ar.Contains(t, s.tags, string(msgExp.Tag), "Tags storage is not initiated for requested tag")
	a.True(t, s.tags[string(msgExp.Tag)].Has(msgExp.ID), "Message.ID is not assigned to tag")
}

func Test_MemoryStorage_MessageSave_Failure_NoID(t *testing.T) {
	s, closer := tsMemoryStorageSetup()
	defer closer()

	ar.EqualError(t, s.MsgSave(&tfMsgAXA_NoID), ErrElementIDNotSet.Error())
	s.usersMu.RLock()
	defer s.usersMu.RUnlock()
	ar.Len(t, s.messages, 0, "unexpected element stored")
}

func Test_MemoryStorage_MessageLoad_Exists(t *testing.T) {
	s, closer := tsMemoryStorageSetup()
	defer closer()

	msgExp := tfMsgAA

	// GIVEN: expected message is in storage
	ar.NoError(t, s.MsgSave(&msgExp))

	msgGot, err := s.MsgLoad(msgExp.ID)
	ar.NoError(t, err)
	a.Equal(t, &msgExp, msgGot, "Message from storage does not match")
}

func Test_MemoryStorage_MessageLoad_NotFound(t *testing.T) {
	s, closer := tsMemoryStorageSetup()
	defer closer()

	// GIVEN: expected message is NOT in storage

	_, err := s.MsgLoad(tfMsgAA.ID)
	ar.Equal(t, ErrElementNotFound, err)
}

// -- section: Tag
func Test_MemoryStorage_TagAddMsgID_First(t *testing.T) {
	s, closer := tsMemoryStorageSetup()
	defer closer()

	tag := Tag("ABC")
	mID := "mID-1"
	s.tagAddMsgID(tag, mID)

	s.tagsMu.RLock()
	defer s.tagsMu.RUnlock()

	tagSet, found := s.tags[string(tag)]
	ar.True(t, found, "no messages associated with tag")

	ar.Equal(t, tagSet.Size(), 1, "mismatch in number of assocaited tags")
	a.True(t, tagSet.Has(mID), "messageID not in set")
}

func Test_MemoryStorage_TagAddMsgID_Next(t *testing.T) {
	s, closer := tsMemoryStorageSetup()
	defer closer()

	tag := Tag("ABC")
	mID1 := "mID-1"
	mID2 := "mID-2"
	s.tagAddMsgID(tag, mID1)
	s.tagAddMsgID(tag, mID2)

	s.tagsMu.RLock()
	tagSet, found := s.tags[string(tag)]
	s.tagsMu.RUnlock()

	ar.True(t, found, "no messages associated with tag")

	ar.Equal(t, tagSet.Size(), 2, "mismatch in number of assocaited tags")
	a.True(t, tagSet.Has(mID1), "messageID not in set")
	a.True(t, tagSet.Has(mID2), "messageID not in set")
}

func Test_MemoryStorage_TagFindMessagesIDs_Exists(t *testing.T) {
	s, closer := tsMemoryStorageSetup()
	defer closer()

	// GIVEN: expected message is in storage
	msgsExp := []Message{tfMsgAA, tfMsgAB, tfMsgBA, tfMsgBB}
	for _, m := range msgsExp {
		mC := m
		ar.NoError(t, s.MsgSave(&mC))
	}

	idsGot, err := s.MsgsIDsFindByTag(tfTagA)
	ar.NoError(t, err)
	ar.Len(t, idsGot, len(msgsExp)-1, "mismatched number of ids returned")
	for _, mExp := range msgsExp {
		if mExp.Tag == tfTagA {
			a.Contains(t, idsGot, mExp.ID, "User from storage does not match")
		}
	}
}

func Test_MemoryStorage_TagFindMessagesIDs_NotFound(t *testing.T) {
	s, closer := tsMemoryStorageSetup()
	defer closer()

	// GIVEN: expected message is in storage
	msgsExp := []Message{tfMsgAA, tfMsgAB, tfMsgBA, tfMsgBB}
	for _, m := range msgsExp {
		mC := m
		ar.NoError(t, s.MsgSave(&mC))
	}

	_, err := s.MsgsIDsFindByTag(tfTagC)
	a.Equal(t, ErrElementNotFound, err)
}

// -- test helpers
func tsMemoryStorageSetup() (*memoryStorage, func()) {
	s := NewMemoryStorage()
	closer := func() {}
	return s, closer
}
