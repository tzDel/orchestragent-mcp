package application

import (
	"context"
	"errors"
	"path/filepath"
	"testing"

	"github.com/tzDel/orchestragent-mcp/internal/domain"
)

const testRepositoryRoot = "/repo/root"

func setupCreateWorktreeUseCase(gitOps *mockGitOperations) (*CreateWorktreeUseCase, *mockSessionRepository) {
	if gitOps == nil {
		gitOps = &mockGitOperations{}
	}

	sessionRepository := newMockSessionRepository()
	useCase := NewCreateWorktreeUseCase(gitOps, sessionRepository, testRepositoryRoot)
	return useCase, sessionRepository
}

func TestCreateWorktreeUseCase_Execute_Success(t *testing.T) {
	// arrange
	createWorktreeUseCase, _ := setupCreateWorktreeUseCase(nil)
	expectedSessionID := "test-session"
	expectedBranchName := "orchestragent-test-session"
	expectedStatus := "open"
	request := CreateWorktreeRequest{SessionID: expectedSessionID}
	ctx := context.Background()

	// act
	response, err := createWorktreeUseCase.Execute(ctx, request)

	// assert
	if err != nil {
		t.Fatalf("Execute() error: %v", err)
	}

	if response.SessionID != expectedSessionID {
		t.Errorf("SessionID = %q, want %q", response.SessionID, expectedSessionID)
	}

	if response.BranchName != expectedBranchName {
		t.Errorf("BranchName = %q, want %q", response.BranchName, expectedBranchName)
	}

	if response.Status != expectedStatus {
		t.Errorf("Status = %q, want %q", response.Status, expectedStatus)
	}
}

func TestCreateWorktreeUseCase_Execute_InvalidSessionID(t *testing.T) {
	// arrange
	createWorktreeUseCase, _ := setupCreateWorktreeUseCase(nil)
	request := CreateWorktreeRequest{SessionID: "Invalid_ID"}
	ctx := context.Background()

	// act
	_, err := createWorktreeUseCase.Execute(ctx, request)

	// assert
	if err == nil {
		t.Error("Execute() expected error for invalid session ID")
	}
}

func TestCreateWorktreeUseCase_Execute_SessionAlreadyExists(t *testing.T) {
	// arrange
	createWorktreeUseCase, sessionRepository := setupCreateWorktreeUseCase(nil)

	sessionID, _ := domain.NewSessionID("test-session")
	session, _ := domain.NewSession(sessionID, "/path")
	sessionRepository.Save(context.Background(), session)

	request := CreateWorktreeRequest{SessionID: "test-session"}
	ctx := context.Background()

	// act
	_, err := createWorktreeUseCase.Execute(ctx, request)

	// assert
	if err == nil {
		t.Error("Execute() expected error for existing session")
	}
}

func TestCreateWorktreeUseCase_Execute_BranchAlreadyExists(t *testing.T) {
	// arrange
	gitOperations := &mockGitOperations{
		branchExistsFunc: func(ctx context.Context, branch string) (bool, error) {
			return true, nil
		},
	}
	createWorktreeUseCase, _ := setupCreateWorktreeUseCase(gitOperations)
	request := CreateWorktreeRequest{SessionID: "test-session"}
	ctx := context.Background()

	// act
	_, err := createWorktreeUseCase.Execute(ctx, request)

	// assert
	if err == nil {
		t.Error("Execute() expected error for existing branch")
	}
}

func TestCreateWorktreeUseCase_Execute_GitOperationFails(t *testing.T) {
	// arrange
	gitOperations := &mockGitOperations{
		createWorktreeFunc: func(ctx context.Context, path string, branch string) error {
			return errors.New("git error")
		},
	}
	createWorktreeUseCase, _ := setupCreateWorktreeUseCase(gitOperations)
	request := CreateWorktreeRequest{SessionID: "test-session"}
	ctx := context.Background()

	// act
	_, err := createWorktreeUseCase.Execute(ctx, request)

	// assert
	if err == nil {
		t.Error("Execute() expected error when git operation fails")
	}
}

func TestCreateWorktreeUseCase_ValidateSessionID_WithValidID_ReturnsSessionID(t *testing.T) {
	// arrange
	createWorktreeUseCase, _ := setupCreateWorktreeUseCase(nil)
	validIDString := "test-session"

	// act
	sessionID, err := createWorktreeUseCase.validateSessionID(validIDString)

	// assert
	if err != nil {
		t.Fatalf("validateSessionID() unexpected error: %v", err)
	}
	if sessionID.String() != validIDString {
		t.Errorf("validateSessionID() returned %q, want %q", sessionID.String(), validIDString)
	}
}

func TestCreateWorktreeUseCase_ValidateSessionID_WithInvalidID_ReturnsError(t *testing.T) {
	// arrange
	createWorktreeUseCase, _ := setupCreateWorktreeUseCase(nil)
	invalidIDString := "Invalid_ID"

	// act
	_, err := createWorktreeUseCase.validateSessionID(invalidIDString)

	// assert
	if err == nil {
		t.Error("validateSessionID() expected error for invalid session ID")
	}
}

func TestCreateWorktreeUseCase_EnsureSessionDoesNotExist_WhenSessionDoesNotExist_ReturnsNoError(t *testing.T) {
	// arrange
	createWorktreeUseCase, _ := setupCreateWorktreeUseCase(nil)
	sessionID, _ := domain.NewSessionID("test-session")
	ctx := context.Background()

	// act
	err := createWorktreeUseCase.ensureSessionDoesNotExist(ctx, sessionID)

	// assert
	if err != nil {
		t.Errorf("ensureSessionDoesNotExist() unexpected error: %v", err)
	}
}

func TestCreateWorktreeUseCase_EnsureSessionDoesNotExist_WhenSessionExists_ReturnsError(t *testing.T) {
	// arrange
	createWorktreeUseCase, sessionRepository := setupCreateWorktreeUseCase(nil)
	sessionID, _ := domain.NewSessionID("test-session")
	session, _ := domain.NewSession(sessionID, "/path")
	sessionRepository.Save(context.Background(), session)
	ctx := context.Background()

	// act
	err := createWorktreeUseCase.ensureSessionDoesNotExist(ctx, sessionID)

	// assert
	if err == nil {
		t.Error("ensureSessionDoesNotExist() expected error when session exists")
	}
}

func TestCreateWorktreeUseCase_EnsureBranchDoesNotExist_WhenBranchDoesNotExist_ReturnsNoError(t *testing.T) {
	// arrange
	gitOperations := &mockGitOperations{
		branchExistsFunc: func(ctx context.Context, branch string) (bool, error) {
			return false, nil
		},
	}
	createWorktreeUseCase, _ := setupCreateWorktreeUseCase(gitOperations)
	branchName := "orchestragent-test-session"
	ctx := context.Background()

	// act
	err := createWorktreeUseCase.ensureBranchDoesNotExist(ctx, branchName)

	// assert
	if err != nil {
		t.Errorf("ensureBranchDoesNotExist() unexpected error: %v", err)
	}
}

func TestCreateWorktreeUseCase_EnsureBranchDoesNotExist_WhenBranchExists_ReturnsError(t *testing.T) {
	// arrange
	gitOperations := &mockGitOperations{
		branchExistsFunc: func(ctx context.Context, branch string) (bool, error) {
			return true, nil
		},
	}
	createWorktreeUseCase, _ := setupCreateWorktreeUseCase(gitOperations)
	branchName := "orchestragent-test-session"
	ctx := context.Background()

	// act
	err := createWorktreeUseCase.ensureBranchDoesNotExist(ctx, branchName)

	// assert
	if err == nil {
		t.Error("ensureBranchDoesNotExist() expected error when branch exists")
	}
}

func TestCreateWorktreeUseCase_BuildWorktreePath_ReturnsCorrectPath(t *testing.T) {
	// arrange
	createWorktreeUseCase, _ := setupCreateWorktreeUseCase(nil)
	sessionID, _ := domain.NewSessionID("test-session")
	expectedWorktreeDirName := "orchestragent-test-session"
	expectedPath := filepath.Join("/repo/root", ".worktrees", expectedWorktreeDirName)

	// act
	worktreePath := createWorktreeUseCase.buildWorktreePath(sessionID)

	// assert
	if worktreePath != expectedPath {
		t.Errorf("buildWorktreePath() returned %q, want %q", worktreePath, expectedPath)
	}
}

func TestCreateWorktreeUseCase_CreateWorktreeAndBranch_Success_ReturnsNoError(t *testing.T) {
	// arrange
	createWorktreeUseCase, _ := setupCreateWorktreeUseCase(nil)
	worktreePath := "/repo/root/.worktrees/orchestragent-test-session"
	branchName := "orchestragent-test-session"
	ctx := context.Background()

	// act
	err := createWorktreeUseCase.createWorktreeAndBranch(ctx, worktreePath, branchName)

	// assert
	if err != nil {
		t.Errorf("createWorktreeAndBranch() unexpected error: %v", err)
	}
}

func TestCreateWorktreeUseCase_CreateWorktreeAndBranch_GitOperationFails_ReturnsError(t *testing.T) {
	// arrange
	gitOperations := &mockGitOperations{
		createWorktreeFunc: func(ctx context.Context, path string, branch string) error {
			return errors.New("git error")
		},
	}
	createWorktreeUseCase, _ := setupCreateWorktreeUseCase(gitOperations)
	worktreePath := "/repo/root/.worktrees/orchestragent-test-session"
	branchName := "orchestragent-test-session"
	ctx := context.Background()

	// act
	err := createWorktreeUseCase.createWorktreeAndBranch(ctx, worktreePath, branchName)

	// assert
	if err == nil {
		t.Error("createWorktreeAndBranch() expected error when git operation fails")
	}
}

func TestCreateWorktreeUseCase_CreateAndSaveSession_Success_ReturnsSession(t *testing.T) {
	// arrange
	createWorktreeUseCase, _ := setupCreateWorktreeUseCase(nil)
	expectedSessionID := "test-session"
	sessionID, _ := domain.NewSessionID(expectedSessionID)
	worktreePath := "/repo/root/.worktrees/orchestragent-test-session"
	ctx := context.Background()

	// act
	session, err := createWorktreeUseCase.createAndSaveSession(ctx, sessionID, worktreePath)

	// assert
	if err != nil {
		t.Fatalf("createAndSaveSession() unexpected error: %v", err)
	}
	if session.ID().String() != expectedSessionID {
		t.Errorf("createAndSaveSession() session ID = %q, want %q", session.ID().String(), expectedSessionID)
	}
	if session.WorktreePath() != worktreePath {
		t.Errorf("createAndSaveSession() worktree path = %q, want %q", session.WorktreePath(), worktreePath)
	}
}

func TestCreateWorktreeUseCase_BuildResponse_ReturnsCorrectResponse(t *testing.T) {
	// arrange
	createWorktreeUseCase, _ := setupCreateWorktreeUseCase(nil)
	expectedSessionID := "test-session"
	expectedBranchName := "orchestragent-test-session"
	expectedStatus := "open"
	sessionID, _ := domain.NewSessionID(expectedSessionID)
	session, _ := domain.NewSession(sessionID, "/repo/root/.worktrees/orchestragent-test-session")

	// act
	response := createWorktreeUseCase.buildResponse(session)

	// assert
	if response.SessionID != expectedSessionID {
		t.Errorf("buildResponse() SessionID = %q, want %q", response.SessionID, expectedSessionID)
	}
	if response.BranchName != expectedBranchName {
		t.Errorf("buildResponse() BranchName = %q, want %q", response.BranchName, expectedBranchName)
	}
	if response.Status != expectedStatus {
		t.Errorf("buildResponse() Status = %q, want %q", response.Status, expectedStatus)
	}
}
