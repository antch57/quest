---
title: "getting started"
weight: 1
---

quest is a lightweight personal cli tool for tracking tasks and discovering live shows — all from your terminal.

## install

**via go install** (recommended if you have go 1.21+):

```bash
go install github.com/antch57/quest@latest
```

**build from source:**

```bash
git clone https://github.com/antch57/quest.git
cd quest
go build -o quest .
```

then move the binary somewhere on your `$PATH`:

```bash
mv quest /usr/local/bin/quest
```

---

## verify the install

```bash
quest --version
quest --help
```

---

## your first tasks

create a todo:

```bash
quest log create --title "try quest"
```

list your open todos:

```bash
quest log list
```

mark it done:

```bash
quest log done --id 1
```

---

## commands

| command | description |
|---|---|
| [`log`](commands/log) | create, track, and complete tasks in your quest log |
| [`jamz`](commands/jamz) | search for upcoming live shows via the jambase api |

---

## environment variables

| variable | required for | description |
|---|---|---|
| `JAMBASE_API_KEY` | `quest jamz` | api key from [jambase.com](https://www.jambase.com/article/jambase-api) |

```bash
export JAMBASE_API_KEY=your_api_key_here
```

{{< callout type="info" >}}
you only need `JAMBASE_API_KEY` if you plan to use `quest jamz`. the `log` command works out of the box with no configuration.
{{< /callout >}}

---

## storage

task data is stored locally at `~/.quest/todos.json`. the directory is created automatically on first use — no database, no account, nothing remote.

---

## what's next

- [log](commands/log) — full guide to managing your quest log
- [jamz](commands/jamz) — finding shows near you