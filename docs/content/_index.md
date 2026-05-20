---
date: '2026-05-19T22:44:12-06:00'
draft: false
type: landing
title: quest
toc: false
---

<div style="text-align: center; padding: 5rem 1rem 4rem;">

<p style="font-size: 0.85rem; letter-spacing: 0.15em; text-transform: uppercase; opacity: 0.5; margin-bottom: 1.25rem;">a personal cli toolkit</p>

<h1 style="font-size: clamp(2.5rem, 6vw, 4rem); font-weight: 800; letter-spacing: -0.04em; line-height: 1.1; margin-bottom: 1.5rem;">
  little tools<br>for everyday quests
</h1>

<p style="font-size: 1.2rem; opacity: 0.7; max-width: 500px; margin: 0 auto 2.5rem; line-height: 1.6;">
  track tasks, discover live shows, and build a workflow that lives entirely in your shell — no accounts, no cloud, no noise.
</p>

<div style="display: flex; gap: 0.75rem; justify-content: center; flex-wrap: wrap;">
  <a href="docs/getting-started" style="display: inline-block; background: #7c3aed; color: white; padding: 0.7rem 1.75rem; border-radius: 0.5rem; font-weight: 600; text-decoration: none; font-size: 0.95rem;">get started →</a>
  <a href="https://github.com/antch57/quest" style="display: inline-block; padding: 0.7rem 1.75rem; border-radius: 0.5rem; font-weight: 600; text-decoration: none; font-size: 0.95rem; border: 1px solid #7c3aed; color: #7c3aed; opacity: 0.75;">view on github</a>
</div>

</div>

---

## what it does

{{< cards >}}
  {{< card link="docs/commands/log" title="quest log" icon="clipboard-list" subtitle="create tasks, set due dates, filter by project, mark things done — all without leaving the terminal." >}}
  {{< card link="docs/commands/jamz" title="quest jamz" icon="music-note" subtitle="search upcoming shows near you via the jambase api. find your next show without touching a browser." >}}
  {{< card title="stays local" icon="lock-closed" subtitle="no accounts. no sync. your data lives in ~/.quest/todos.json — plain json, yours forever." >}}
{{< /cards >}}

---

## install

```bash
go install github.com/antch57/quest@latest
```

or build from source:

```bash
git clone https://github.com/antch57/quest.git
cd quest && go build -o quest . && mv quest /usr/local/bin/
```

---

## a taste

{{< asciinema
    file="demo.cast"
    theme="dracula"
    speed="1"
    autoplay="true"
    loop="true"
    font-size="16px"
>}}

<!-- ```bash
# add a task with a due date and project
quest log create -t "ship the feature" -d 05-30-2026 -p "work"

# see what's open
quest log list

# filter by project
quest log list -p "work"

# find a show near you
quest jamz search --city denver

# check it off
quest log done --id 1
``` -->

---

<div style="text-align: center; padding: 3rem 1rem 2rem;">
  <p style="font-size: 1.05rem; opacity: 0.65; margin-bottom: 1.5rem;">ready to begin your quest?</p>
  <a href="docs/getting-started" style="display: inline-block; background: #7c3aed; color: white; padding: 0.7rem 1.75rem; border-radius: 0.5rem; font-weight: 600; text-decoration: none; font-size: 0.95rem;">read the docs →</a>
</div>
