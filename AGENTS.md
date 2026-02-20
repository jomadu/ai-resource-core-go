# AGENTS.md

## Work Tracking System

Tasks are tracked in `TASK.md` at repository root.

Task format:
```markdown
## TASK-001
- Priority: 1-5 (1=highest)
- Status: TODO/IN_PROGRESS/BLOCKED/DONE
- Dependencies: [TASK-XXX, ...]
- Description: Task description
```

Manual editing. Tasks auto-increment. Keep all tasks (including DONE) in file.

## Quick Reference

- Edit `TASK.md` - Manage tasks
- `go test ./...` - Run tests (when initialized)
- `go build ./...` - Build packages (when initialized)

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

Excludes: `testdata/`, `.git/`

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

Last verified: 2026-02-20

**Working:**
- `bd` command available for task tracking
- Specification structure in `specs/` directory
- TASK.md defines implementation plan

**Not working:**
- Go module not initialized (no go.mod)
- Implementation directories not created
- PLAN.md not created

**Rationale:**
- Project in bootstrap phase
- Specifications complete, implementation pending