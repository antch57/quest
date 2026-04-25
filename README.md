# quest

my personal CLI hub: starting with notes and todos, growing over time into a daily command center.

## wizard's field notes

### v0.1.0 - Fork in the Enchanted Road

the party reached a crossroads and chose a clearer trail.

progress made in this chapter:

- moved todo operations under the `quest log` command domain
- reorganized command files into a future-ready structure under `cmd/log`
- switched to explicit flags for key actions:
  - `quest log create --title "..."`
  - `quest log done --id <id>`
  - `quest log delete --id <id>`
  - `quest log edit --id <id> --title "..."`
- kept `quest log nuke` as a dedicated high-impact action
- rewired `main.go` so the `log` domain is now the primary quest board

next quests on the map:

- add command-layer tests for `quest log` flows (`create`, `done`, `delete`, `edit`, `list`)
- sharpen help text so required flags are obvious for every command
- keep refining `quest log` UX before locking the final interface style
- open the gates to web-powered utility commands once `quest log` is fully battle-tested

### v0.0.1 - Apprentice's First Parchment
quest started as a simple todo cli, but the longer-term idea is to turn it into a personal cli hub: one place for lightweight daily tools, quick logging, and small web-powered commands.

the first big shift will be moving the current todo workflow under a dedicated `log` subcommand. instead of top-level commands like `quest create`, the log flow would become flag-driven and live under a single surface, for example:

```sh
quest log --create "buy milk" --due 01-01-2030
quest log --done 3
quest log --list
quest log --delete 4
```

that keeps the current task system intact while making room for quest to grow into something broader.

from there, the project can branch into a few directions:

- keep evolving `quest log` into a better capture and review tool for tasks, notes, and daily activity
- add web-powered commands that reach out to external services for useful daily info
- experiment with commands for weather, headlines, github summaries, links, or other personal utility flows
- keep using the project as a place to practice more go: http requests, json decoding, config handling, and larger cli structure

the goal is not to make a giant framework. the goal is to build a small wizardy command center that feels personal, useful, and fun to extend.

## trail ahead

this is the forward-looking roadmap for quest.

### phase 1 - harden the quest log

- add command-layer tests for create, done, delete, edit, and list
- tighten usage and help output for required flags like --title and --id
- smooth rough edges in command UX and error messaging

### phase 2 - improve daily use

- add quick note capture flow (lightweight text entry)
- add optional tags or categories for better organization
- add better filtering and summary views for daily review

### phase 3 - open the outer world

- introduce first web-powered command (for example weather or headlines)
- add simple config support for external service settings
- keep each new command small, practical, and easy to maintain

### phase 4 - evolve into a personal hub

- keep expanding with tools that support real day-to-day workflow
- maintain a clean command layout as new domains are added
- continue using quest as a go learning lab while building useful features