# quest

![Go](https://img.shields.io/badge/Go-1.26-blue?logo=go)
![CI](https://github.com/antch57/quest/actions/workflows/ci.yml/badge.svg)
![Release](https://github.com/antch57/quest/actions/workflows/release-please.yml/badge.svg)

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
./quest log create --title "Buy milk" --due 2026-05-01
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
- Improve help output and flag validation for all commands
- Add support for user-defined projects/tags on todos
- Polish CLI UX (error messages, confirmation prompts, etc.)

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
