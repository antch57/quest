[![Latest Release](https://img.shields.io/github/v/release/antch57/quest)](https://github.com/antch57/quest/releases)
[![Build Status](https://img.shields.io/github/actions/workflow/status/antch57/quest/ci.yml?branch=main)](https://github.com/antch57/quest/actions)
[![Go Reference](https://pkg.go.dev/badge/github.com/antch57/quest.svg)](https://pkg.go.dev/github.com/antch57/quest)
[![Go Report Card](https://goreportcard.com/badge/github.com/antch57/quest)](https://goreportcard.com/report/github.com/antch57/quest)
[![Test Coverage](https://codecov.io/github/antch57/quest/graph/badge.svg?token=7ZEEPESVHP)](https://codecov.io/github/antch57/quest)

# quest

Your personal command-line HQ for todos, notes, and daily tools—crafted to grow into a playful, powerful, and endlessly extensible adventurer’s toolkit.

---

## About

quest is a personal CLI hub: starting with notes and todos, growing into a daily command center for lightweight tools and quick logging. The goal: a small, personal, and fun-to-extend command wizard.

## Quick Start

### Install

Clone and build from source:

```sh
git clone https://github.com/antch57/quest.git
cd quest
go build -o quest .
```

Or, if you have Go installed:

```sh
go install github.com/antch57/quest@latest
```

### Usage

Add a todo:
```sh
./quest log create --title "Buy milk" --due 05-01-2026
```

List todos:
```sh
./quest log list
```

Mark a todo as done:
```sh
./quest log done --id 1
```

Delete a todo:
```sh
./quest log delete --id 1
```

See all commands and flags:
```sh
./quest log --help
```

## Roadmap (planned features & improvements)

### 2026-04-25 (current plan)
- Add command-layer tests for all `quest log` subcommands (create, list, edit, done, delete, nuke)
- Add a `quest weather` command to fetch and display current weather for a given location
- ~~Improve help output and flag validation for all commands~~ (done in v0.1.1)
- ~~Add support for user-defined projects/tags on todos~~ (done in v0.1.1)
- ~~Polish CLI UX (error messages, confirmation prompts, etc.)~~ (done in v0.1.1)

### Future ideas
- Add a quick note capture flow (`quest note`)
- Add a `quest headlines` command for daily news
- Add config file support for API keys and user preferences
- Add recurring tasks and reminders
- Add web-powered commands for GitHub issues, calendar, or links
- Optional: Web dashboard or TUI for richer review

---

## Release Notes

Release notes and changelog are managed automatically by [release-please](https://github.com/google-github-actions/release-please-action) and can be found in the [GitHub Releases](../../releases) tab.

## How to contribute

1. Fork and clone the repo
2. Create a feature branch
3. Use Conventional Commits (e.g., `feat:`, `fix:`, `chore:`)
4. Open a PR (template auto-populates)
5. All changes require PR review and CI to pass
6. Run `gofmt`, `go vet`, and tests before pushing
7. See the PR template for checklist

## Roadmap update process

- Add new roadmap entries at the top with the current date and keep them forward-looking
- Never delete old entries—this is a living project plan

---
