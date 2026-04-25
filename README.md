# quest

simple little todo cli tool to practice go.

## wizard's field notes

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