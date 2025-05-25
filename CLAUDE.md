# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

memov2 is a Go CLI/TUI application for managing markdown-based memos and todos. It uses file-based storage (no database) with markdown files containing YAML frontmatter metadata.

## Build & Test Commands

```bash
make build      # Build binary (output: ./memov2)
make install    # Install to $GOPATH/bin
make test       # Run all tests with coverage (go test ./... -cover)
go test ./internal/domain/ -run TestMemoFile -v  # Run a single test
```

## Architecture

Layered architecture with dependency injection:

```
cmd/             → Cobra CLI commands (memos, todos, config subcommands)
internal/app/    → App container (DI: wires config, services, repos, platform)
internal/service/    → Business logic (memo, todo, config services)
internal/repositories/  → Data access (file system operations for memo, todo, weekly)
internal/platform/   → External integrations (editor, filesystem, trash)
internal/domain/     → Entities (MemoFile, TodoFile, WeeklyFile, HeadingBlock)
internal/interfaces/ → Centralized interface definitions for all layers
internal/search/     → Shared search logic (romaji conversion, SKK dictionary, memo search)
internal/ui/tui/     → Bubbletea TUI (browse + search modes, tab to switch)
internal/config/     → TOML-based configuration
internal/common/     → Structured AppError with ErrorType, slog-based logging
```

Entry point: `main.go` → `cmd.Execute()` → `cmd/root.go` → subcommands.

DI flow: `cmd/app/app.go:InitializeApp()` → `internal/app/app.go:NewApp()` → wires config, services, repos, platform.

## Key Conventions

- **Interfaces**: All defined centrally in `internal/interfaces/` (service.go, repository.go, domain.go, platform.go, config.go). Implementations live in their respective packages.
- **Error handling**: Use `common.New(ErrorType, message)` and `common.Wrap(err, ErrorType, message)` from `internal/common/errors.go`. ErrorTypes: config, filesystem, validation, repository, service, ui.
- **File naming**: Memo/todo files follow `YYYYMMDDDAY000000_TYPE_title.md` pattern.
- **Testing**: Tests use testify assertions. Repos and services have `test_helpers.go` files and `testdata/` directories with `*.golden` files for expected output.
- **Mocks**: Repository mocks in `internal/repositories/mock/`.
- **Display width**: Use `runewidth.StringWidth()` (from `go-runewidth`) for terminal column width of strings containing Japanese characters. Do not use `len()` or `%-*s` for alignment — they count bytes, not display columns.

## Spec-Driven Development

This project uses Kiro-style spec-driven development (see `AGENTS.md`). Feature specs live in `.kiro/specs/`, project steering docs in `.kiro/steering/` (product.md, tech.md, structure.md).

## Language

Development guidelines in AGENTS.md specify: think in English, write spec documents in the language configured in spec.json.
