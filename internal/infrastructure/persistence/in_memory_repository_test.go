package persistence

import (
	"context"
	"testing"

	"github.com/tzDel/orchestrAIgent/internal/domain"
)

func TestInMemoryRepository_SaveAndFind(t *testing.T) {
	// arrange
	repository := NewInMemorySessionRepository()
	ctx := context.Background()
	sessionID, _ := domain.NewSessionID("test-session")
	session, _ := domain.NewSession(sessionID, "/path/to/worktree")

	// act
	err := repository.Save(ctx, session)
	if err != nil {
		t.Fatalf("Save() error: %v", err)
	}

	found, err := repository.FindByID(ctx, sessionID)

	// assert
	if err != nil {
		t.Fatalf("FindByID() error: %v", err)
	}

	if found.ID().String() != session.ID().String() {
		t.Errorf("FindByID() returned wrong session")
	}
}

func TestInMemoryRepository_FindByID_NotFound(t *testing.T) {
	// arrange
	repository := NewInMemorySessionRepository()
	ctx := context.Background()
	sessionID, _ := domain.NewSessionID("nonexistent")

	// act
	_, err := repository.FindByID(ctx, sessionID)

	// assert
	if err == nil {
		t.Error("FindByID() expected error for non-existent session")
	}
}

func TestInMemoryRepository_Exists(t *testing.T) {
	// arrange
	repository := NewInMemorySessionRepository()
	ctx := context.Background()
	sessionID, _ := domain.NewSessionID("test-session")

	// act
	exists, err := repository.Exists(ctx, sessionID)

	// assert
	if err != nil {
		t.Fatalf("Exists() error: %v", err)
	}
	if exists {
		t.Error("Exists() returned true for non-existent session")
	}

	// arrange
	session, _ := domain.NewSession(sessionID, "/path")
	repository.Save(ctx, session)

	// act
	exists, err = repository.Exists(ctx, sessionID)

	// assert
	if err != nil {
		t.Fatalf("Exists() error: %v", err)
	}
	if !exists {
		t.Error("Exists() returned false for existing session")
	}
}

func TestInMemoryRepository_FindAll_EmptyRepository(t *testing.T) {
	// arrange
	repository := NewInMemorySessionRepository()
	ctx := context.Background()

	// act
	sessions, err := repository.FindAll(ctx)

	// assert
	if err != nil {
		t.Fatalf("FindAll() error: %v", err)
	}
	if len(sessions) != 0 {
		t.Errorf("FindAll() returned %d sessions, want 0", len(sessions))
	}
}

func TestInMemoryRepository_FindAll_MultipleSessions(t *testing.T) {
	// arrange
	repository := NewInMemorySessionRepository()
	ctx := context.Background()

	sessionID1, _ := domain.NewSessionID("session-one")
	session1, _ := domain.NewSession(sessionID1, "/path/one")

	sessionID2, _ := domain.NewSessionID("session-two")
	session2, _ := domain.NewSession(sessionID2, "/path/two")

	sessionID3, _ := domain.NewSessionID("session-three")
	session3, _ := domain.NewSession(sessionID3, "/path/three")

	repository.Save(ctx, session1)
	repository.Save(ctx, session2)
	repository.Save(ctx, session3)

	// act
	sessions, err := repository.FindAll(ctx)

	// assert
	if err != nil {
		t.Fatalf("FindAll() error: %v", err)
	}
	if len(sessions) != 3 {
		t.Errorf("FindAll() returned %d sessions, want 3", len(sessions))
	}

	sessionIDs := make(map[string]bool)
	for _, session := range sessions {
		sessionIDs[session.ID().String()] = true
	}

	if !sessionIDs["session-one"] || !sessionIDs["session-two"] || !sessionIDs["session-three"] {
		t.Error("FindAll() did not return all expected sessions")
	}
}
