package application

import (
	"context"

	"github.com/tzDel/orchestrAIgent/internal/domain"
)

type GetSessionsRequest struct {
}

type SessionDTO struct {
	SessionID    string `json:"sessionId"`
	WorktreePath string `json:"worktreePath"`
	BranchName   string `json:"branchName"`
	Status       string `json:"status"`
	LinesAdded   int    `json:"linesAdded"`
	LinesRemoved int    `json:"linesRemoved"`
}

type GetSessionsResponse struct {
	Sessions []SessionDTO `json:"sessions"`
}

type GetSessionsUseCase struct {
	gitOperations     domain.GitOperations
	sessionRepository domain.SessionRepository
	baseBranch        string
}

func NewGetSessionsUseCase(
	gitOperations domain.GitOperations,
	sessionRepository domain.SessionRepository,
	baseBranch string,
) *GetSessionsUseCase {
	return &GetSessionsUseCase{
		gitOperations:     gitOperations,
		sessionRepository: sessionRepository,
		baseBranch:        baseBranch,
	}
}

func (useCase *GetSessionsUseCase) Execute(ctx context.Context, request GetSessionsRequest) (*GetSessionsResponse, error) {
	sessions, err := useCase.sessionRepository.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	sessionDTOs := make([]SessionDTO, 0, len(sessions))
	for _, session := range sessions {
		diffStats, err := useCase.gitOperations.GetDiffStats(ctx, session.WorktreePath(), useCase.baseBranch)
		if err != nil {
			// Continue with zero stats on error
			diffStats = &domain.GitDiffStats{LinesAdded: 0, LinesRemoved: 0}
		}

		dto := useCase.buildSessionDTO(session, diffStats)
		sessionDTOs = append(sessionDTOs, dto)
	}

	return &GetSessionsResponse{
		Sessions: sessionDTOs,
	}, nil
}

func (useCase *GetSessionsUseCase) buildSessionDTO(session *domain.Session, diffStats *domain.GitDiffStats) SessionDTO {
	return SessionDTO{
		SessionID:    session.ID().String(),
		WorktreePath: session.WorktreePath(),
		BranchName:   session.BranchName(),
		Status:       string(session.Status()),
		LinesAdded:   diffStats.LinesAdded,
		LinesRemoved: diffStats.LinesRemoved,
	}
}
