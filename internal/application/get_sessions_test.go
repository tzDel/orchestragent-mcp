package application

import (
	"context"
	"testing"

	"github.com/tzDel/orchestragent-mcp/internal/domain"
)

func TestGetSessionsUseCase_NoSessions(t *testing.T) {
	// arrange
	mockGitOps := &MockGitOperations{}
	mockRepo := &MockSessionRepository{
		sessions: make(map[string]*domain.Session),
	}
	useCase := NewGetSessionsUseCase(mockGitOps, mockRepo, "main")
	ctx := context.Background()
	request := GetSessionsRequest{}

	// act
	response, err := useCase.Execute(ctx, request)

	// assert
	if err != nil {
		t.Fatalf("Execute() error: %v", err)
	}
	if len(response.Sessions) != 0 {
		t.Errorf("Execute() returned %d sessions, want 0", len(response.Sessions))
	}
}

func TestGetSessionsUseCase_MultipleSessions(t *testing.T) {
	// arrange
	sessionID1, _ := domain.NewSessionID("session-one")
	session1, _ := domain.NewSession(sessionID1, "/path/session-one")

	sessionID2, _ := domain.NewSessionID("session-two")
	session2, _ := domain.NewSession(sessionID2, "/path/session-two")

	mockGitOps := &MockGitOperations{
		diffStats: map[string]*domain.GitDiffStats{
			"session-one": {
				LinesAdded:   42,
				LinesRemoved: 17,
			},
			"session-two": {
				LinesAdded:   5,
				LinesRemoved: 3,
			},
		},
	}

	mockRepo := &MockSessionRepository{
		sessions: map[string]*domain.Session{
			"session-one": session1,
			"session-two": session2,
		},
	}

	useCase := NewGetSessionsUseCase(mockGitOps, mockRepo, "main")
	ctx := context.Background()
	request := GetSessionsRequest{}

	// act
	response, err := useCase.Execute(ctx, request)

	// assert
	if err != nil {
		t.Fatalf("Execute() error: %v", err)
	}
	if len(response.Sessions) != 2 {
		t.Errorf("Execute() returned %d sessions, want 2", len(response.Sessions))
	}

	sessionMap := make(map[string]SessionDTO)
	for _, session := range response.Sessions {
		sessionMap[session.SessionID] = session
	}

	session1DTO := sessionMap["session-one"]
	if session1DTO.LinesAdded != 42 {
		t.Errorf("Session one LinesAdded = %d, want 42", session1DTO.LinesAdded)
	}
	if session1DTO.LinesRemoved != 17 {
		t.Errorf("Session one LinesRemoved = %d, want 17", session1DTO.LinesRemoved)
	}
	if session1DTO.Status != "open" {
		t.Errorf("Session one Status = %s, want open", session1DTO.Status)
	}

	session2DTO := sessionMap["session-two"]
	if session2DTO.LinesAdded != 5 {
		t.Errorf("Session two LinesAdded = %d, want 5", session2DTO.LinesAdded)
	}
	if session2DTO.LinesRemoved != 3 {
		t.Errorf("Session two LinesRemoved = %d, want 3", session2DTO.LinesRemoved)
	}
	if session2DTO.Status != "open" {
		t.Errorf("Session two Status = %s, want open", session2DTO.Status)
	}
}

func TestGetSessionsUseCase_GitStatsError_ContinuesWithOtherSessions(t *testing.T) {
	// arrange
	sessionID1, _ := domain.NewSessionID("session-one")
	session1, _ := domain.NewSession(sessionID1, "/path/session-one")

	sessionID2, _ := domain.NewSessionID("session-two")
	session2, _ := domain.NewSession(sessionID2, "/path/session-two")

	mockGitOps := &MockGitOperations{
		diffStats: map[string]*domain.GitDiffStats{
			"session-two": {
				LinesAdded:   5,
				LinesRemoved: 3,
			},
		},
		shouldFailForSession: "session-one",
	}

	mockRepo := &MockSessionRepository{
		sessions: map[string]*domain.Session{
			"session-one": session1,
			"session-two": session2,
		},
	}

	useCase := NewGetSessionsUseCase(mockGitOps, mockRepo, "main")
	ctx := context.Background()
	request := GetSessionsRequest{}

	// act
	response, err := useCase.Execute(ctx, request)

	// assert
	if err != nil {
		t.Fatalf("Execute() error: %v", err)
	}

	if len(response.Sessions) != 2 {
		t.Errorf("Execute() returned %d sessions, want 2", len(response.Sessions))
	}

	sessionMap := make(map[string]SessionDTO)
	for _, session := range response.Sessions {
		sessionMap[session.SessionID] = session
	}

	session1DTO := sessionMap["session-one"]
	if session1DTO.LinesAdded != 0 || session1DTO.LinesRemoved != 0 {
		t.Error("Session one should have zero stats when git operation fails")
	}

	session2DTO := sessionMap["session-two"]
	if session2DTO.LinesAdded != 5 {
		t.Errorf("Session two LinesAdded = %d, want 5", session2DTO.LinesAdded)
	}
}
