[![Latest Release](https://img.shields.io/github/v/release/antch57/quest)](https://github.com/antch57/quest/releases)
[![docs](https://img.shields.io/badge/docs-purple?logo=hugo&logoColor=black)](https://antch57.github.io/quest)
[![ci](https://github.com/antch57/quest/actions/workflows/ci.yml/badge.svg)](https://github.com/antch57/quest/actions/workflows/ci.yml)
[![Test Coverage](https://codecov.io/github/antch57/quest/graph/badge.svg?token=7ZEEPESVHP)](https://codecov.io/github/antch57/quest)
[![Go Reference](https://pkg.go.dev/badge/github.com/antch57/quest.svg)](https://pkg.go.dev/github.com/antch57/quest)
[![Go Report Card](https://goreportcard.com/badge/github.com/antch57/quest)](https://goreportcard.com/report/github.com/antch57/quest)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

# quest

a personal cli for todos, notes, really whatever you need it to be. quest is lightweight and built to grow into an endlessly extensible toolkit for your daily grind.

---

## quick start

### install

clone and build from source:

```sh
git clone https://github.com/antch57/quest.git
cd quest
go build -o quest .
```

or, if you have go installed:

```sh
go install github.com/antch57/quest@latest
```

## subcommands

### log

track and manage tasks from the terminal.

#### usage

```sh
quest log create --title "buy milk" --due 05-01-2026
quest log list
quest log update --id 1 --title "buy oat milk"
quest log --help
```

### jamz

find upcoming shows using the jambase api.

#### setup

set your jambase api key before running searches:

```sh
export JAMBASE_API_KEY="your-api-key"
```

#### usage

```sh
quest jamz search --city denver
quest jamz search --artist goose --date 2026-06-01
quest jamz search --city denver --radius 50 --limit 10
quest jamz --help
```

## release notes

release notes and changelog are managed automatically by [release-please](https://github.com/google-github-actions/release-please-action) and can be found in the [github releases](../../releases) tab.

## how to contribute

1. fork or clone repo
1. create a feature branch
1. use conventional commits (e.g., `feat:`, `fix:`, `chore:`)
1. open a pr
1. see the pr template for checklist
1. all changes require pr review and ci to pass

## license

MIT. see [LICENSE](LICENSE) for details.
