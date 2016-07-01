package main

import (
	"errors"
	"sync"

	"github.com/fatih/set"
)

var (
	ErrElementIDNotSet = errors.New("Storage: element ID not set")
	ErrElementNotFound = errors.New("Storage: element not found")
)

// memoryStorage provides in memory storage for users
// All functions are thread safe.
type memoryStorage struct {
	// users is a storage for a users.
	// Keyed by User.ID
	users map[string]*User
	// usersMu is RW mutex protecting users map
	usersMu sync.RWMutex

	// messages is a storage for messages.
	// Keyed by Message.ID.
	messages map[string]*Message
	// messagesMu is RW mutex protecting messages map.
	messagesMu sync.RWMutex

	// tags keeps association between messages and tags
	// Keyed by tag with sets of message.ID as value.
	tags map[string]*set.Set
	// tagsMu is RW mutex protecting tags map.
	tagsMu sync.RWMutex
}

// NewMemoryStorage returns empty memory storage
func NewMemoryStorage() *memoryStorage {
	return &memoryStorage{
		users:    make(map[string]*User),
		messages: make(map[string]*Message),
		tags:     make(map[string]*set.Set),
	}
}

// UserSave persists single user.
// ErrElementIDNotSet error is returned if user ID is not set.
func (s *memoryStorage) UserSave(u *User) error {
	if u.ID == "" {
		return ErrElementIDNotSet
	}
	s.usersMu.Lock()
	defer s.usersMu.Unlock()
	s.users[u.ID] = u

	return nil
}

// UserLoad retrieves single user from storage by ID.
// ErrElementNotFound is returned if element could not be found.
func (s *memoryStorage) UserLoad(id string) (*User, error) {
	s.usersMu.RLock()
	defer s.usersMu.RUnlock()
	u, found := s.users[id]
	if !found {
		return nil, ErrElementNotFound
	}
	return u, nil
}

// UserFindByName retrieves single user entity from storage by its Name.
// ErrElementNotFound is returned if user could not be found.
// TODO: optimise me -> search is implemented as naive O(N) scan.
func (s *memoryStorage) UserFindByName(name string) (*User, error) {
	s.usersMu.RLock()
	defer s.usersMu.RUnlock()
	for _, o := range s.users {
		if o.Name == name {
			return o, nil
		}
	}
	return nil, ErrElementNotFound
}

// MsgSave persists single message.
// Error ErrElementIDNotSet is dispatched when message ID is not set.
func (s *memoryStorage) MsgSave(m *Message) error {
	if m.ID == "" {
		return ErrElementIDNotSet
	}
	s.messagesMu.Lock()
	defer s.messagesMu.Unlock()
	s.messages[m.ID] = m

	s.tagAddMsgID(m.Tag, m.ID)

	return nil
}

// MsgLoad retrieves single message from storage by ID.
// ErrElementNotFound is returned if message could not be found.
func (s *memoryStorage) MsgLoad(id string) (*Message, error) {
	s.messagesMu.RLock()
	defer s.messagesMu.RUnlock()
	m, found := s.messages[id]
	if !found {
		return nil, ErrElementNotFound
	}
	return m, nil
}

// tagAddMsgID is a helper which adds messageID to a tag
func (s *memoryStorage) tagAddMsgID(tag Tag, mID string) {
	s.tagsMu.Lock()
	defer s.tagsMu.Unlock()

	ts, found := s.tags[string(tag)]
	if found {
		ts.Add(mID)
		return
	}

	s.tags[string(tag)] = set.New(mID)
}

// MsgsIDsFindByTag returns list of ids of messages associated with given tag.
// ErrElementNotFound is returned if tag is unknown (no message is associated)
func (s *memoryStorage) MsgsIDsFindByTag(tag Tag) ([]string, error) {
	s.tagsMu.RLock()
	defer s.tagsMu.RUnlock()

	ts, found := s.tags[string(tag)]
	if !found {
		return []string{}, ErrElementNotFound
	}

	out := make([]string, 0, ts.Size())
	ts.Each(func(item interface{}) bool {
		out = append(out, item.(string))
		return true
	})

	return out, nil
}
