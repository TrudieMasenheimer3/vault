package vault

import (
	"testing"
)

func TestOIDCLoginDeduplication(t *testing.T) {
	store := NewIdentityStore()
	
	group := &Group{
		ID:              "group-1",
		Name:            "oidc-group",
		MemberEntityIds: []string{},
	}
	
	err := store.SaveGroup(group)
	if err != nil {
		t.Fatalf("failed to save group: %v", err)
	}

	entityID := "user-entity-1"
	groupIDs := []string{"group-1"}

	// Trigger login 5 times sequentially
	for i := 0; i < 5; i++ {
		err := store.OIDCLogin(entityID, groupIDs)
		if err != nil {
			t.Fatalf("login failed at iteration %d: %v", i, err)
		}
	}

	// Load group and assert member count
	loadedGroup, err := store.LoadGroup("group-1")
	if err != nil {
		t.Fatalf("failed to load group: %v", err)
	}

	if len(loadedGroup.MemberEntityIds) != 1 {
		t.Errorf("expected exactly 1 member entity ID, got %d: %v", len(loadedGroup.MemberEntityIds), loadedGroup.MemberEntityIds)
	}

	if loadedGroup.MemberEntityIds[0] != entityID {
		t.Errorf("expected member entity ID to be %q, got %q", entityID, loadedGroup.MemberEntityIds[0])
	}
}

func TestDeduplicateOnLoadSaveUpdate(t *testing.T) {
	store := NewIdentityStore()
	
	group := &Group{
		ID:              "group-2",
		Name:            "test-group",
		MemberEntityIds: []string{"user-1", "user-1", "user-2", "user-1"},
	}

	// Test deduplication on Save
	err := store.SaveGroup(group)
	if err != nil {
		t.Fatalf("failed to save group: %v", err)
	}

	loadedGroup, err := store.LoadGroup("group-2")
	if err != nil {
		t.Fatalf("failed to load group: %v", err)
	}

	expected := []string{"user-1", "user-2"}
	if len(loadedGroup.MemberEntityIds) != len(expected) {
		t.Fatalf("expected %d members, got %d", len(expected), len(loadedGroup.MemberEntityIds))
	}
	for i, v := range expected {
		if loadedGroup.MemberEntityIds[i] != v {
			t.Errorf("expected member at index %d to be %q, got %q", i, v, loadedGroup.MemberEntityIds[i])
		}
	}

	// Test deduplication on Update
	loadedGroup.MemberEntityIds = append(loadedGroup.MemberEntityIds, "user-2", "user-3", "user-3")
	err = store.UpdateGroup(loadedGroup)
	if err != nil {
		t.Fatalf("failed to update group: %v", err)
	}

	loadedGroup, err = store.LoadGroup("group-2")
	if err != nil {
		t.Fatalf("failed to load group: %v", err)
	}

	expectedUpdate := []string{"user-1", "user-2", "user-3"}
	if len(loadedGroup.MemberEntityIds) != len(expectedUpdate) {
		t.Fatalf("expected %d members, got %d", len(expectedUpdate), len(loadedGroup.MemberEntityIds))
	}
	for i, v := range expectedUpdate {
		if loadedGroup.MemberEntityIds[i] != v {
			t.Errorf("expected member at index %d to be %q, got %q", i, v, loadedGroup.MemberEntityIds[i])
		}
	}
}
