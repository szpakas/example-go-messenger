package main

// tmMemoryStorageMock is a wrapped MemoryStorage with tracking/custom error capabilities for testing
type tmMemoryStorageMock struct {
	*memoryStorage

	inUserSaveCalled bool
	outUserSaveErr   error

	inUserLoadCalled bool
	outUserLoadErr   error

	inUserFindCalled bool
	outUserFindErr   error

	inMsgSaveCalled bool
	outMsgSaveErr   error

	inMsgLoadCalled bool
	outMsgLoadErr   error

	inMsgFindCalled bool
	outMsgFindErr   error
}

func (s *tmMemoryStorageMock) UserSave(u *User) error {
	s.inUserSaveCalled = true

	if s.outUserSaveErr != nil {
		return s.outUserSaveErr
	}
	return s.memoryStorage.UserSave(u)
}

func (s *tmMemoryStorageMock) UserLoad(id string) (*User, error) {
	s.inUserLoadCalled = true

	if s.outUserLoadErr != nil {
		return nil, s.outUserLoadErr
	}
	return s.memoryStorage.UserLoad(id)
}

func (s *tmMemoryStorageMock) UserFindByName(name string) (*User, error) {
	s.inUserFindCalled = true

	if s.outUserFindErr != nil {
		return nil, s.outUserFindErr
	}
	return s.memoryStorage.UserFindByName(name)
}

func (s *tmMemoryStorageMock) MsgSave(m *Message) error {
	s.inMsgSaveCalled = true

	if s.outMsgSaveErr != nil {
		return s.outMsgSaveErr
	}
	return s.memoryStorage.MsgSave(m)
}

func (s *tmMemoryStorageMock) MsgLoad(id string) (*Message, error) {
	s.inMsgLoadCalled = true

	if s.outMsgLoadErr != nil {
		return nil, s.outMsgLoadErr
	}
	return s.memoryStorage.MsgLoad(id)
}

func (s *tmMemoryStorageMock) MsgsIDsFindByTag(tag Tag) ([]string, error) {
	s.inMsgFindCalled = true

	if s.outMsgFindErr != nil {
		return []string{}, s.outMsgFindErr
	}
	return s.memoryStorage.MsgsIDsFindByTag(tag)
}

func NewTmMemoryStorageMock() *tmMemoryStorageMock {
	sto := NewMemoryStorage()
	return &tmMemoryStorageMock{
		memoryStorage: sto,
	}
}
