package mcp

import (
	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/tzDel/orchestragent-mcp/internal/application"
)

type CreateWorktreeArgs struct {
	SessionID string `json:"sessionId" jsonschema:"required" jsonschema_description:"The unique identifier for the session"`
}

type CreateWorktreeOutput struct {
	SessionID    string `json:"sessionId"`
	WorktreePath string `json:"worktreePath"`
	BranchName   string `json:"branchName"`
	Status       string `json:"status"`
}

type RemoveSessionArgs struct {
	SessionID string `json:"sessionId" jsonschema:"required" jsonschema_description:"Session identifier"`
	Force     bool   `json:"force" jsonschema_description:"Skip safety checks and force removal"`
}

type RemoveSessionOutput struct {
	SessionID          string `json:"sessionId"`
	RemovedAt          string `json:"removedAt,omitempty"`
	HasUnmergedChanges bool   `json:"hasUnmergedChanges"`
	UnmergedCommits    int    `json:"unmergedCommits"`
	UncommittedFiles   int    `json:"uncommittedFiles"`
	Warning            string `json:"warning,omitempty"`
}

type GetSessionsArgs struct {
}

type GetSessionsOutput struct {
	Sessions []SessionOutput `json:"sessions"`
}

type SessionOutput struct {
	SessionID    string `json:"sessionId"`
	WorktreePath string `json:"worktreePath"`
	BranchName   string `json:"branchName"`
	Status       string `json:"status"`
	LinesAdded   int    `json:"linesAdded"`
	LinesRemoved int    `json:"linesRemoved"`
}

type MCPServer struct {
	mcpServer             *mcpsdk.Server
	createWorktreeUseCase *application.CreateWorktreeUseCase
	removeSessionUseCase  *application.RemoveSessionUseCase
	getSessionsUseCase    *application.GetSessionsUseCase
}
