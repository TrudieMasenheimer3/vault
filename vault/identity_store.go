package vault

import (
	"errors"
	"sync"
)

type IdentityStore struct {
	mu     sync.RWMutex
	groups map[string]*Group
}

func NewIdentityStore() *IdentityStore {
	return &IdentityStore{
		groups: make(map[string]*Group),
	}
}

func (s *IdentityStore) LoadGroup(id string) (*Group, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	g, exists := s.groups[id]
	if !exists {
		return nil, errors.New("group not found")
	}
	g.MemberEntityIds = Deduplicate(g.MemberEntityIds)
	return g, nil
}

func (s *IdentityStore) SaveGroup(g *Group) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if g == nil {
		return errors.New("nil group")
	}
	g.MemberEntityIds = Deduplicate(g.MemberEntityIds)
	s.groups[g.ID] = g
	return nil
}

func (s *IdentityStore) UpdateGroup(g *Group) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if g == nil {
		return errors.New("nil group")
	}
	g.MemberEntityIds = Deduplicate(g.MemberEntityIds)
	s.groups[g.ID] = g
	return nil
}

func (s *IdentityStore) AddMemberToGroup(groupID string, entityID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	g, exists := s.groups[groupID]
	if !exists {
		return errors.New("group not found")
	}
	
	g.MemberEntityIds = Deduplicate(g.MemberEntityIds)

	existsMember := false
	for _, id := range g.MemberEntityIds {
		if id == entityID {
			existsMember = true
			break
		}
	}
	if !existsMember {
		g.MemberEntityIds = append(g.MemberEntityIds, entityID)
	}
	return nil
}

func (s *IdentityStore) OIDCLogin(entityID string, groupIDs []string) error {
	for _, groupID := range groupIDs {
		if err := s.AddMemberToGroup(groupID, entityID); err != nil {
			return err
		}
	}
	return nil
}
