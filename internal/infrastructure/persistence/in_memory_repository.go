package persistence

import (
	"context"
	"fmt"
	"sync"

	"github.com/tzDel/orchestragent-mcp/internal/domain"
)

type InMemorySessionRepository struct {
	mutex    sync.RWMutex
	sessions map[string]*domain.Session
}

func NewInMemorySessionRepository() *InMemorySessionRepository {
	return &InMemorySessionRepository{
		sessions: make(map[string]*domain.Session),
	}
}

func (repository *InMemorySessionRepository) Save(ctx context.Context, session *domain.Session) error {
	repository.mutex.Lock()
	defer repository.mutex.Unlock()

	repository.sessions[session.ID().String()] = session
	return nil
}

func (repository *InMemorySessionRepository) FindByID(ctx context.Context, sessionID domain.SessionID) (*domain.Session, error) {
	repository.mutex.RLock()
	defer repository.mutex.RUnlock()

	session, exists := repository.sessions[sessionID.String()]
	if !exists {
		return nil, fmt.Errorf("session not found: %s", sessionID.String())
	}

	return session, nil
}

func (repository *InMemorySessionRepository) FindAll(ctx context.Context) ([]*domain.Session, error) {
	repository.mutex.RLock()
	defer repository.mutex.RUnlock()

	sessions := make([]*domain.Session, 0, len(repository.sessions))
	for _, session := range repository.sessions {
		sessions = append(sessions, session)
	}

	return sessions, nil
}

func (repository *InMemorySessionRepository) Exists(ctx context.Context, sessionID domain.SessionID) (bool, error) {
	repository.mutex.RLock()
	defer repository.mutex.RUnlock()

	_, exists := repository.sessions[sessionID.String()]
	return exists, nil
}

func (repository *InMemorySessionRepository) Delete(ctx context.Context, sessionID domain.SessionID) error {
	repository.mutex.Lock()
	defer repository.mutex.Unlock()

	if _, exists := repository.sessions[sessionID.String()]; !exists {
		return fmt.Errorf("session not found: %s", sessionID.String())
	}

	delete(repository.sessions, sessionID.String())
	return nil
}
