package application

import (
	"context"
	"errors"
	"path/filepath"

	"github.com/tzDel/orchestragent-mcp/internal/domain"
)

type mockGitOperations struct {
	createWorktreeFunc        func(ctx context.Context, path string, branch string) error
	removeWorktreeFunc        func(ctx context.Context, path string, force bool) error
	branchExistsFunc          func(ctx context.Context, branch string) (bool, error)
	hasUncommittedChangesFunc func(ctx context.Context, worktreePath string) (bool, int, error)
	hasUnpushedCommitsFunc    func(ctx context.Context, baseBranch string, sessionBranch string) (int, error)
	deleteBranchFunc          func(ctx context.Context, branchName string, force bool) error
	getDiffStatsFunc          func(ctx context.Context, worktreePath string, baseBranch string) (*domain.GitDiffStats, error)
}

type MockGitOperations struct {
	diffStats            map[string]*domain.GitDiffStats
	shouldFailForSession string
}

func (mock *mockGitOperations) CreateWorktree(ctx context.Context, path string, branch string) error {
	if mock.createWorktreeFunc != nil {
		return mock.createWorktreeFunc(ctx, path, branch)
	}
	return nil
}

func (mock *mockGitOperations) RemoveWorktree(ctx context.Context, path string, force bool) error {
	if mock.removeWorktreeFunc != nil {
		return mock.removeWorktreeFunc(ctx, path, force)
	}
	return nil
}

func (mock *mockGitOperations) BranchExists(ctx context.Context, branch string) (bool, error) {
	if mock.branchExistsFunc != nil {
		return mock.branchExistsFunc(ctx, branch)
	}
	return false, nil
}

func (mock *mockGitOperations) HasUncommittedChanges(ctx context.Context, worktreePath string) (bool, int, error) {
	if mock.hasUncommittedChangesFunc != nil {
		return mock.hasUncommittedChangesFunc(ctx, worktreePath)
	}
	return false, 0, nil
}

func (mock *mockGitOperations) HasUnpushedCommits(ctx context.Context, baseBranch string, sessionBranch string) (int, error) {
	if mock.hasUnpushedCommitsFunc != nil {
		return mock.hasUnpushedCommitsFunc(ctx, baseBranch, sessionBranch)
	}
	return 0, nil
}

func (mock *mockGitOperations) DeleteBranch(ctx context.Context, branchName string, force bool) error {
	if mock.deleteBranchFunc != nil {
		return mock.deleteBranchFunc(ctx, branchName, force)
	}
	return nil
}

func (mock *mockGitOperations) GetDiffStats(ctx context.Context, worktreePath string, baseBranch string) (*domain.GitDiffStats, error) {
	if mock.getDiffStatsFunc != nil {
		return mock.getDiffStatsFunc(ctx, worktreePath, baseBranch)
	}
	return &domain.GitDiffStats{LinesAdded: 0, LinesRemoved: 0}, nil
}

func (mock *MockGitOperations) CreateWorktree(ctx context.Context, path string, branch string) error {
	return nil
}

func (mock *MockGitOperations) RemoveWorktree(ctx context.Context, path string, force bool) error {
	return nil
}

func (mock *MockGitOperations) BranchExists(ctx context.Context, branch string) (bool, error) {
	return false, nil
}

func (mock *MockGitOperations) HasUncommittedChanges(ctx context.Context, worktreePath string) (bool, int, error) {
	return false, 0, nil
}

func (mock *MockGitOperations) HasUnpushedCommits(ctx context.Context, baseBranch string, sessionBranch string) (int, error) {
	return 0, nil
}

func (mock *MockGitOperations) DeleteBranch(ctx context.Context, branchName string, force bool) error {
	return nil
}

func (mock *MockGitOperations) GetDiffStats(ctx context.Context, worktreePath string, baseBranch string) (*domain.GitDiffStats, error) {
	sessionID := filepath.Base(worktreePath)

	if mock.shouldFailForSession != "" {
		if sessionID == mock.shouldFailForSession {
			return nil, errors.New("git diff failed")
		}
	}

	if stats, exists := mock.diffStats[sessionID]; exists {
		return stats, nil
	}

	return &domain.GitDiffStats{LinesAdded: 0, LinesRemoved: 0}, nil
}

type mockSessionRepository struct {
	sessions map[string]*domain.Session
}

func newMockSessionRepository() *mockSessionRepository {
	return &mockSessionRepository{
		sessions: make(map[string]*domain.Session),
	}
}

func (mock *mockSessionRepository) Save(ctx context.Context, session *domain.Session) error {
	mock.sessions[session.ID().String()] = session
	return nil
}

func (mock *mockSessionRepository) FindByID(ctx context.Context, sessionID domain.SessionID) (*domain.Session, error) {
	session, exists := mock.sessions[sessionID.String()]
	if !exists {
		return nil, errors.New("not found")
	}
	return session, nil
}

func (mock *mockSessionRepository) FindAll(ctx context.Context) ([]*domain.Session, error) {
	sessions := make([]*domain.Session, 0, len(mock.sessions))
	for _, session := range mock.sessions {
		sessions = append(sessions, session)
	}
	return sessions, nil
}

func (mock *mockSessionRepository) Exists(ctx context.Context, sessionID domain.SessionID) (bool, error) {
	_, exists := mock.sessions[sessionID.String()]
	return exists, nil
}

func (mock *mockSessionRepository) Delete(ctx context.Context, sessionID domain.SessionID) error {
	if _, exists := mock.sessions[sessionID.String()]; !exists {
		return errors.New("not found")
	}
	delete(mock.sessions, sessionID.String())
	return nil
}

type MockSessionRepository struct {
	sessions map[string]*domain.Session
}

func (mock *MockSessionRepository) Save(ctx context.Context, session *domain.Session) error {
	mock.sessions[session.ID().String()] = session
	return nil
}

func (mock *MockSessionRepository) FindByID(ctx context.Context, sessionID domain.SessionID) (*domain.Session, error) {
	session, exists := mock.sessions[sessionID.String()]
	if !exists {
		return nil, errors.New("not found")
	}
	return session, nil
}

func (mock *MockSessionRepository) FindAll(ctx context.Context) ([]*domain.Session, error) {
	sessions := make([]*domain.Session, 0, len(mock.sessions))
	for _, session := range mock.sessions {
		sessions = append(sessions, session)
	}
	return sessions, nil
}

func (mock *MockSessionRepository) Exists(ctx context.Context, sessionID domain.SessionID) (bool, error) {
	_, exists := mock.sessions[sessionID.String()]
	return exists, nil
}

func (mock *MockSessionRepository) Delete(ctx context.Context, sessionID domain.SessionID) error {
	if _, exists := mock.sessions[sessionID.String()]; !exists {
		return errors.New("not found")
	}
	delete(mock.sessions, sessionID.String())
	return nil
}
