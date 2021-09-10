# airtable-dnd
A CLI to manage a D&amp;D campaign from the command line with Airtable.

Currently *only partially implemented* (WIP) and only able to record encounters
as they are completed alongside an optional ability to display the initiative
order of characters and enemies.

## Build

`go build .`

## Usage

This script has a number of commands and will eventually feature more.

### Manage Items

`./airtable-dnd items <item_json_file_paths...>` uses
a [5etools](https://5etools.com/) JSON file or set of JSON files containing item
data to populate weight and value data in a party inventory table.

### Manage Encounters

`./airtable-dnd encounter <encounter_json_file_path> [submit/slots]` where the
path to an encounter JSON file points to a file with the following example
format:

```
{
  "Encounter": {
    "Name": "s2_l1_r2",
    "Session": 2,
    "Level": 1,
    "Room": 2,
    "XP": 30
  },
  "Participants": [
    {"Name": "Bandit", "Initiative": 2, "HP": 12},
    {"Name": "Burglar", "Initiative": 2, "HP": 23},
    {"Name": "Mook", "Initiative": 2}
  ]
}
```

## Current Features

  - Encounter submission/recordkeeping
  - Encounter initiative order
  - Management of party loot

## Planned Features

- Random generation of entities (NPCs, etc.)
- Management of characters

## Related Projects

  - [Saber](https://github.com/alexSafatli/saber) which is intended to eventually be an underlying engine to this CLI
