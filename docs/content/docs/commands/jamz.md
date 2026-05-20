---
title: "jamz"
---

find upcoming live shows near you — powered by the [jambase api](https://data.jambase.com/api/docs/getting-started)

---

## overview

`quest jamz` lets you search for shows from the terminal. results are pulled from jambase and displayed in a clean table with venue, date, and door time.

| subcommand | description |
|---|---|
| `search` | search for upcoming shows |

### requirements

jamz requires a jambase api key set as an environment variable:

```bash
export JAMBASE_API_KEY=your_api_key_here
```

{{< callout type="info" >}}
you can get a free api key at [data.jambase.com](https://data.jambase.com/api/docs/getting-started).
{{< /callout >}}

---

## search

search for upcoming shows. all flags are optional — running it with just `--city` is the most common way to use it.

```bash
quest jamz search [options]
```

### flags

| flag | alias | default | description |
|---|---|---|---|
| `--city` | `-c` | | city to search around |
| `--country` | | `US` | two-letter iso country code |
| `--artist` | `-a` | | filter results by artist name |
| `--venue` | `-v` | | filter results by venue name |
| `--date` | `-d` | | search for shows on a specific date (`yyyy-mm-dd`) |
| `--radius` | `-r` | `25` | search radius in miles |
| `--limit` | `-n` | `25` | max number of results to return (max: 50) |

### examples

```bash
# shows near a city
quest jamz search --city denver

# expand the search radius
quest jamz search --city denver --radius 50

# filter by artist
quest jamz search --city chicago -a "arctic monkeys"

# filter by venue
quest jamz search --city nashville -v "ryman auditorium"

# shows on a specific date
quest jamz search --city austin -d 2026-06-15

# search in another country
quest jamz search --city london --country GB

# limit results
quest jamz search --city seattle -n 10
```

---

## output

results are displayed as a table in your terminal:

```
 #  | name                        | date         | doors  | venue               | address
----+-----------------------------+--------------+--------+---------------------+------------------
 1  | goose                       | 2026-06-01   | 7:00pm | red rocks           | morrison, co
 2  | phish                       | 2026-06-03   | 6:30pm | ogden theatre       | denver, co
```

{{< callout type="warning" >}}
the `--limit` flag caps at 50. values above 50 will return an error.
{{< /callout >}}