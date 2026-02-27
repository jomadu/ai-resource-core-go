# AGENTS.md

## Work Tracking System

Tasks are tracked in `TODO.json` at repository root.

Task format:
```json
{
  "id": "TASK-001",
  "priority": 1,
  "status": "TODO",
  "dependencies": ["TASK-XXX"],
  "description": "Task description",
  "comments": []
}
```

Fields:
- `id`: Task identifier (TASK-XXX)
- `priority`: 1-5 (1=highest)
- `status`: TODO/IN_PROGRESS/BLOCKED/DONE
- `dependencies`: Array of task IDs
- `description`: Task description
- `comments`: Array of implementation notes (agents append as needed)

Manual editing. Tasks auto-increment. Keep all tasks (including DONE) in file.

## Feature Input

`TASK.md` contains feature requirements and specifications for rooda.

## Quick Reference

- Edit `TODO.json` - Manage tasks
- `jq '.tasks[] | select(.id == "TASK-XXX")' TODO.json` - Show single task
- `jq '.tasks[] | select(.status != "DONE")' TODO.json` - List incomplete tasks
- `jq '.tasks[] | select(.status == "TODO" and (.dependencies | length == 0 or all(. as $dep | any($dep == .tasks[].id and .tasks[].status == "DONE"))))' TODO.json` - List ready tasks
- `make help` - Show all available commands
- `make test` - Run all tests (auto-initializes submodule)
- `make test-conformance` - Run conformance tests only
- `make build` - Build all packages
- `make lint` - Run linters
- `make update-spec` - Update spec to latest version

## Planning System

`PLAN.md` documents the current plan (not yet created).

## Build/Test/Lint Commands

Go project (not yet initialized):
- `go mod init` - Initialize module
- `go test ./...` - Run tests
- `go build ./...` - Build packages
- `go vet ./...` - Lint code

## Specification Definition

Specifications live in `specs/`. A spec template exists at `specs/TEMPLATE.md`.

Format: Markdown with structured sections (Purpose, Requirements, Implementation Notes, Examples).

## Implementation Definition

Location: `pkg/airesource/` (public API), `internal/` (private implementation)

Patterns:
- `pkg/airesource/*.go` - Public API surface
- `internal/schema/*.go` - Schema validation
- `internal/template/*.go` - Mustache rendering
- `internal/types/*.go` - Type validation

Excludes: `.git/`

## Audit Output

Audit results written to `AUDIT.md` at repository root.

## Quality Criteria

**Specifications:**
- All requirements testable
- Examples provided
- Implementation notes clear

**Implementation:**
- Passes `go test ./...`
- Passes `go vet ./...`
- Public API minimal and documented
- Error types structured

**Refactoring triggers:**
- Spec/implementation divergence
- Test failures
- Unclear error messages

## Operational Learnings

Last verified: 2026-02-25

**Working:**
- Specification structure in `specs/` directory
- TASK.md defines implementation plan
- Makefile is the standard interface for all development tasks
- Conformance test submodule is REQUIRED (not optional)

**Not working:**
- Go module not initialized (no go.mod)
- Implementation directories not created
- PLAN.md not created

**Rationale:**
- Project in bootstrap phase
- Specifications complete, implementation pending

**Troubleshooting:**
- If tests fail with submodule error, run `make test` (auto-initializes submodule)
- Submodule must be initialized before running conformance tests
- Use `make update-spec` to update spec to latest version