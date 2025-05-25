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

```bash
# Create a new memo
memov2 memos new "Meeting Notes"

# Create a memo with a category
memov2 memos new "Design Notes" --category work/projects

# List all categories (useful for piping to fzf/peco)
memov2 memos categories

# Create a memo with an interactively selected category
memov2 memos new "Minutes" --category "$(memov2 memos categories | fzf)"
memov2 memos new "Minutes" --category "$(memov2 memos categories | peco)"
# If fzf/peco is cancelled (empty selection), the memo is created without a category

# Browse and search memos in the TUI
memov2 memos browse

# List all memos
memov2 memos list

# Open a memo
memov2 memos open "memos/work/20250114Mon150405_memo_notes.md"

# Search memos with romaji support (for use with fzf)
memov2 memos search "kaigi"

# Search with match context (shows where each match was found)
memov2 memos search --context "kaigi"

# Rename a memo (updates both title and filename)
memov2 memos rename "work/20250114Mon150405_memo_notes.md" "New Title"
memov2 memos rename "work/20250114Mon150405_memo_notes.md"  # prompts interactively

# Generate a weekly report
memov2 memos weekly

# Generate an index file of all memos
memov2 memos index
```

### Tasks

```bash
# Create today's task file
memov2 todos new

# Recreate today's task file from scratch
memov2 todos new --truncate

# Generate a weekly task report
memov2 todos weekly
```

### Config

```bash
# Show current configuration
memov2 config show

# Open the config file in the configured editor
memov2 config edit
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

Set `editor` to the executable name, not a shell alias — `exec.Command` bypasses the shell, so aliases defined in `.zshrc` / `.bashrc` are not resolved.

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

Filename: `YYYYMMDDDAY000000_memo_title.md`

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

## TUI (`memov2 memos browse`)

The TUI has two modes — Browse and Search — switchable with `Tab`.

### Browse mode

Displays memos in a category tree.

| Key | Action |
|-----|--------|
| `j` / `k` | Move up / down |
| `Ctrl+u` / `Ctrl+d` | Jump 10 lines |
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

## Interactive search with fzf

`memov2 memos search` supports romaji-to-Japanese conversion (SKK dictionary based). Combined with fzf, it enables interactive memo search without launching the TUI.

### Basic (title and path only)

```bash
memov2 memos list | fzf \
  --disabled \
  --bind "change:reload:memov2 memos search {q}" \
| cut -f2 | xargs memov2 memos open
```

**How it works:**

1. `memov2 memos list` shows the initial memo list
2. As you type, `memov2 memos search {q}` is called and the list updates in real time
3. Romaji input (e.g. `kaigi`) matches Japanese memos (e.g. those containing「会議」)
4. The selected memo's path is extracted with `cut -f2` and opened with `memov2 memos open`

### With match context (shows where each match was found)

Use the `--context` (`-c`) flag to display the match type and matched content alongside each result.

```bash
memov2 memos list | fzf \
  --disabled \
  --delimiter=$'\t' \
  --with-nth='1,3,4' \
  --bind "change:reload:memov2 memos search --context {q}" \
| cut -f2 | xargs memov2 memos open
```

**Output format (tab-separated, 4 fields):**

```
Title	path/to/file.md	[Content]	## Design > matched line from body
Title	path/to/file.md	[Title]  	Title
Title	path/to/file.md	[Heading]	Heading text
```

- `--with-nth='1,3,4'` displays title, match type, and matched content (path is hidden)
- `cut -f2` extracts the path and passes it to `memos open`

### Rename a memo with fzf

```bash
memov2 memos list | fzf | cut -f2 | xargs memov2 memos rename
```

Select a memo in fzf, then enter a new title interactively.

### Shell function examples

```bash
# Add to .bashrc / .zshrc
memo-search() {
  memov2 memos list | fzf \
    --disabled \
    --bind "change:reload:memov2 memos search {q}" \
  | cut -f2 | xargs memov2 memos open
}

# Search with match context
memo-search-context() {
  memov2 memos list | fzf \
    --disabled \
    --delimiter=$'\t' \
    --with-nth='1,3,4' \
    --bind "change:reload:memov2 memos search --context {q}" \
  | cut -f2 | xargs memov2 memos open
}

memo-rename() {
  memov2 memos list | fzf | cut -f2 | xargs memov2 memos rename
}

# Create a memo with an interactively selected category (cancelled fzf = no category)
memo-new() {
  local category
  category=$(memov2 memos categories | fzf --prompt="Category: ")
  memov2 memos new "$1" --category "$category"
}
```

## Architecture

```
cmd/                 → Cobra CLI command definitions
internal/
  app/               → DI container (wires config, services, repositories)
  service/           → Business logic (memo, todo, config)
  repositories/      → Data access (filesystem operations)
  platform/          → External integrations (editor, filesystem, trash)
  domain/            → Entities (MemoFile, TodoFile, WeeklyFile)
  interfaces/        → Centralized interface definitions for all layers
  ui/tui/            → Bubbletea TUI (browse + search modes)
  config/            → TOML-based configuration management
  common/            → Error handling, logging
```

## Development

```bash
make build    # build binary
make install  # install to $GOPATH/bin
make test     # run all tests with coverage
```

## License

MIT
