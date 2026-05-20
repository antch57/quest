---
title: "log"
weight: 2
---

manage your quest log — create, track, and complete tasks from the terminal

---

## overview

`quest log` is the core command for managing your personal task list. todos are stored locally in `~/.quest/todos.json` so everything stays on your machine.

| subcommand | description |
|---|---|
| `create` | add a new todo |
| `list` | view todos with optional filters |
| `edit` | update a todo's fields |
| `done` | mark a todo as complete |
| `delete` | soft-delete a todo by id |
| `nuke` | wipe the entire todo store |

---

## create

add a new todo to your quest log.

```bash
quest log create --title "your task title"
```

### flags

| flag | alias | required | description |
|---|---|---|---|
| `--title` | `-t` | yes | title of the todo |
| `--due` | `-d` | no | due date in `mm-dd-yyyy` format |
| `--project` | `-p` | no | project or folder label |

### examples

```bash
# minimal — just a title
quest log create --title "fix login bug"

# with a due date
quest log create -t "submit report" -d 05-30-2026

# scoped to a project
quest log create -t "write tests" -p "quest"

# all options
quest log create -t "deploy v2" -d 06-01-2026 -p "infra"
```

---

## list

display todos from your quest log. by default only incomplete, non-deleted todos are shown.

```bash
quest log list [options]
```

### flags

| flag | alias | description |
|---|---|---|
| `--all` | `-a` | show all todos, including completed ones |
| `--done` | `-d` | show only completed todos |
| `--today` | | todos created today |
| `--week` | | todos created in the last 7 days |
| `--month` | | todos created this month |
| `--overdue` | | todos past their due date |
| `--project` | `-p` | filter by project name |

### examples

```bash
# default view — open todos only
quest log list

# everything, including done
quest log list --all

# only completed tasks
quest log list --done

# due today or earlier, not yet done
quest log list --overdue

# filter by project
quest log list -p "quest"

# this week's todos for a specific project
quest log list --week -p "infra"
```

---

## edit

update one or more fields on an existing todo.

```bash
quest log edit --id <id> [options]
```

### flags

| flag | alias | required | description |
|---|---|---|---|
| `--id` | | yes | id of the todo to edit |
| `--title` | `-t` | no | new title |
| `--due` | `-d` | no | new due date (`mm-dd-yyyy`) |
| `--project` | `-p` | no | new project label |
| `--clear-due` | | no | remove the due date |
| `--done` | | no | mark as complete |
| `--undone` | | no | mark as incomplete |

{{< callout type="info" >}}
`--done` and `--undone` cannot be used together. `--due` and `--clear-due` cannot be used together.
{{< /callout >}}

### examples

```bash
# rename a task
quest log edit --id 3 -t "refactor auth module"

# set a due date
quest log edit --id 3 -d 05-25-2026

# move to a different project
quest log edit --id 3 -p "backend"

# clear the due date
quest log edit --id 3 --clear-due

# mark complete via edit
quest log edit --id 3 --done
```

---

## done

mark a todo as complete by its id.

```bash
quest log done --id <id>
```

### flags

| flag | required | description |
|---|---|---|
| `--id` | yes | id of the todo to mark done |

### examples

```bash
quest log done --id 5
```

---

## delete

soft-delete a todo by id. the record is kept in the store but hidden from all list views.

```bash
quest log delete --id <id>
```

### flags

| flag | required | description |
|---|---|---|
| `--id` | yes | id of the todo to delete |

### examples

```bash
quest log delete --id 2
```

---

## nuke

permanently removes `~/.quest/todos.json`. this cannot be undone.

```bash
quest log nuke
```

you will be prompted to confirm before anything is deleted.

```
are you sure you want to nuke all tasks? (y/n):
```

{{< callout type="warning" >}}
`nuke` deletes the entire store file. all todos — including completed and deleted ones — are gone permanently.
{{< /callout >}}

---

## storage

todos are persisted as json at `~/.quest/todos.json`. the directory is created automatically on first use.

each todo has the following shape:

```json
{
  "id": "1",
  "title": "fix login bug",
  "done": false,
  "deleted": false,
  "created_at": "2026-05-20T10:00:00Z",
  "due_date": "2026-05-30T00:00:00Z",
  "project": "quest"
}
```