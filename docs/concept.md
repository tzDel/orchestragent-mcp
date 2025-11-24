# Concept

## Why

AI coding agents (Copilot, Claude, GPT, Gemini) typically work directly in your main repository, making it hard to:
- **Isolate** multiple agents working simultaneously
- **Control** filesystem access for compliance requirements
- **Parallelize** work without conflicts
- **Manage** agent lifecycles systematically

This MCP server solves these problems by giving each agent an isolated git worktree with restricted filesystem access.

## Vision

**Current (MVP):** Works with existing MCP-compatible CLI tools (Copilot CLI, Claude Code CLI). Agent calls `create_worktree()` → works in isolated worktree → developer reviews and merges manually.

**Future:** Server spawns and monitors agent processes directly. Developer triggers agent via CLI/IDE → server creates worktree + spawns agent process → agent works autonomously → developer reviews/merges via UI.

---

## Capabilities

**MVP:**
- Create isolated git worktree + branch per agent
- Track agent state (worktree path, branch, status)
- Cleanup worktrees/branches
- Provide git info (diffs, commits, files) for review

**Future:**
- Spawn/monitor/terminate agent processes
- IDE plugin integration (VSCode/IntelliJ)
- CLI tool (`orchestraigent` command)
- Persistence, logging, multi-repo support

---

## Tech Stack

**Language:** Go (single binary, cross-platform, excellent process management)

**Core:**
- Go stdlib for process execution (`os/exec`), filesystem ops, JSON
- Git CLI via shell commands (`git worktree`, `git merge`, etc.)
- MCP SDK (`github.com/modelcontextprotocol/go-sdk`)

**Configuration:** YAML/JSON files (repository path, base branch, test command)

**Future:** SQLite/file persistence, CLI tool (`cobra`), structured logging