# memov2

A Markdown-based memo and task management CLI/TUI tool.

All data is stored directly as Markdown files on the filesystem — no database required. Files can be edited with any editor and are well-suited for version control with Git.

## Features

- **File-based**: Stored as Markdown files with YAML frontmatter for metadata
- **Category tree**: Organize memos in a hierarchical category structure
- **Timestamp naming**: Files are automatically organized by date and time
- **TUI browser**: Browse and search memos interactively in the terminal
- **Weekly reports**: Automatically generate weekly summaries for memos and tasks

## Installation

```bash
git clone https://github.com/hirotoni/memov2.git
cd memov2
make build    # produces ./memov2
make install  # installs to $GOPATH/bin
```

Requires Go 1.24.3 or later.

## Usage

### Memos

`new`, `search`, and `rename` are interactive: they open an embedded terminal UI
(no external tools required). See [Interactive commands](#interactive-commands) for
the keybindings.

```bash
# Create a memo: pick a category in a TUI, type a title
memov2 memos new

# Search memos and open the selection (romaji-aware incremental search)
memov2 memos search

# Non-interactive search (pass a query): prints "title<TAB>path" per matched memo
memov2 memos search 会議
memov2 memos search 会議 --json   # machine-readable output

# Rename a memo: pick it in a TUI, type a new title
memov2 memos rename

# Browse memos in a category tree (TUI)
memov2 memos browse

# List all memos (output: "title<TAB>path", absolute paths)
memov2 memos list
memov2 memos list --json    # machine-readable output

# List with relative paths instead of absolute paths
memov2 memos list --short   # -s

# List all categories
memov2 memos categories
memov2 memos categories --json

# Show a single memo by path (raw markdown by default)
memov2 memos show "work/20250114Mon150405_memo_notes.md"
memov2 memos show "work/20250114Mon150405_memo_notes.md" --json

# Open a memo by path
memov2 memos open "memos/work/20250114Mon150405_memo_notes.md"

# Generate a weekly report (writes memos/weekly_report.md, then opens it)
memov2 memos weekly

# Generate an index file of all memos (writes memos/index.md, then opens it)
memov2 memos index
```

Both `weekly` and `index` run a tidy pass first (see [Tidy behavior](#tidy-behavior-weekly--index)), and overwrite their output file (`weekly_report.md` / `index.md`) on each run.

### Tasks

```bash
# Create today's task file
memov2 todos new

# Recreate today's task file from scratch
memov2 todos new --truncate   # -t

# Generate a weekly task report (writes todos/weekly_report.md, then opens it)
memov2 todos weekly
```

Unlike `memos weekly`, `todos weekly` does not run the tidy pass; it only reads the task files for the period and overwrites `todos/weekly_report.md`.

### Config

```bash
# Show current configuration
memov2 config show
memov2 config show --json   # machine-readable output

# Open the config file in the configured editor
memov2 config edit
```

## AI / automation

The read-only data commands support a `--json` flag that writes a single valid
JSON value to **stdout only** (logs and errors stay on stderr), so output can be
piped straight to `jq` or consumed by an AI agent or script:

| Command | Without `--json` | With `--json` |
| --- | --- | --- |
| `memos list` | `title<TAB>path` per line | array of `{ title, path, category, date }` |
| `memos categories` | `a/b/c` per line | array of `{ category, path }` |
| `memos show <path>` | raw markdown of the file | `{ title, path, category, date, top_level_content, headings[] }` |
| `memos search <query>` | `title<TAB>path` per matched memo | array of `{ title, path, category, matches[] }` |
| `config show` | structured log lines | `{ base_dir, todos_dir, memos_dir, todos_days_to_seek }` |

Notes:

- These commands are strictly read-only and non-interactive (they never open a TUI or prompt for input).
- `memos search` with **no** query still launches the interactive TUI; pass a query argument for non-interactive output.
- `memos show` resolves its path under the memos directory; a path that escapes it (e.g. `../../etc/passwd`) is rejected with a non-zero exit and nothing on stdout.
- `<path>` is relative to the memos directory (as printed by `memos list --short`) or absolute.

```bash
# Examples
memov2 memos list --json | jq '.[].title'
memov2 memos search 会議 --json | jq '.[] | {title, matches: (.matches | length)}'
memov2 memos show "$(memov2 memos list --short --json | jq -r '.[0].path')" --json | jq .top_level_content
```

## Configuration

`~/.config/memov2/config.toml` is auto-generated on first run.

```toml
base_dir = "~/.config/memov2/dailymemo/"   # base directory for all files
memos_foldername = "memos/"                 # subdirectory for memos
todos_foldername = "todos/"                 # subdirectory for tasks
todos_daystoseek = 10                       # how many days back to inherit tasks
editor = "vi"                               # editor executable
editor_args = ["{path}"]                    # arguments passed to the editor (template)
```

### Editor configuration

Set `editor` to the executable name, not a shell alias — the editor is launched without a shell, so aliases defined in `.zshrc` / `.bashrc` are not resolved.

`editor_args` is a TOML array where each element is passed as a separate argument. Two template variables are available: `{path}` (file path) and `{basedir}` (base directory).

**Examples by editor:**

```toml
# vi / vim
editor = "vi"
editor_args = ["{path}"]

# Neovim + Neo-tree (open basedir as root, reveal the file)
editor = "nvim"
editor_args = ["-c", "cd {basedir} | Neotree reveal", "{path}"]

# VS Code (open as workspace, jump to file)
editor = "code"
editor_args = ["--folder-uri", "{basedir}", "--goto", "{path}:7"]
```

## File formats

### Memo file

Filename: `YYYYMMDDDAYHHMMSS_memo_title.md` (e.g. `20250114Mon150405_memo_notes.md`)

```markdown
---
category: ["work", "projects"]
---

# Meeting Notes

## Topic 1

Content...

## Topic 2

Content...
```

### Task file

Filename: `YYYYMMDDDAY_todos.md`

```markdown
# 20250214Fri

## meetings

## todos

- [ ] Task A
- [ ] Task B
  - [ ] Subtask 1
  - [ ] Subtask 2
- [x] Completed task

## wanttodos
```

## Tidy behavior (`weekly` / `index`)

Before building their output, `memos weekly` and `memos index` run a tidy pass over the memos directory. The tidy pass:

1. Reads each memo's `category` frontmatter and moves the file to the matching subdirectory under `memos/` (e.g. `category: ["work", "projects"]` → `memos/work/projects/`). The frontmatter is the source of truth; the file's current location is corrected to match it.
2. Removes directories left empty by the moves.

`weekly_report.md` and `index.md` are excluded from the move. If you edit a memo's category frontmatter by hand, the next `weekly` or `index` run is what relocates the file on disk.

## TUI (`memov2 memos browse`)

The TUI has two modes — Browse and Search — switchable with `Tab`.

### Browse mode

Displays memos in a category tree.

| Key | Action |
|-----|--------|
| `j` / `k` | Move up / down |
| `Ctrl+u` / `Ctrl+d` | Jump 10 items |
| `l` | Expand directory / open file |
| `h` | Collapse directory |
| `>` / `<` | Expand all / collapse all |
| `p` | Toggle preview panel |
| `N` | Create new memo in selected category |
| `r` | Rename memo |
| `d` | Duplicate memo |
| `D` | Delete memo (moves to trash) |
| `c` | Change category |
| `Tab` | Switch to Search mode |
| `q` | Quit |

### Search mode

Fuzzy search across titles, categories, body text, and headings. Supports romaji-to-Japanese conversion for Japanese memo search.

## Interactive commands

`memos search`, `memos new`, and `memos rename` open self-contained terminal UIs
(no external tools required). They share a common look and keybindings.

| Command | Flow |
|---------|------|
| `memos search` | Type to filter (romaji-aware), `Ctrl+n` / `Ctrl+p` (or `↓` / `↑`) to move the highlight, `Enter` opens the highlighted memo in the editor |
| `memos new` | Pick a category, "no category", or type a new category path (a `+ new category "…"` row appears) — then type a title; the memo is created and opened |
| `memos rename` | Pick a memo, then type a new title; both the filename and the in-file title are updated |

Common keys: type to filter, `Ctrl+n` / `Ctrl+p` (or arrows) to navigate, `Enter` to
select/confirm, `Esc` to cancel (in the title step, `Esc` goes back to the list).

Search matches titles, categories, body text, and headings, with romaji-to-Japanese
conversion (SKK dictionary based) so input like `kaigi` matches memos containing「会議」.

## Development

```bash
make build    # build binary
make install  # install to $GOPATH/bin
make test     # run all tests with coverage
```

## License

MIT
