# pmt - Prompt Manager Tool

A CLI tool for managing AI prompt snippets. Save, organize, and quickly apply your commonly used prompts for GitHub Copilot and other AI assistants.

Think of it as `git stash` for your AI prompts.

## Features

- Save prompts with automatic git project detection
- Interactive selection UI with fuzzy search
- Automatic clipboard copying
- Organize by type and tags
- Filter and search capabilities
- Simple YAML-based storage

## Installation

### From Source

```bash
git clone https://github.com/sunny/pmt.git
cd pmt
go build -o pmt
sudo mv pmt /usr/local/bin/
```

Or install directly:

```bash
go install github.com/sunny/pmt@latest
```

## Quick Start

```bash
# Save a prompt
pmt push "Fix memory leak in async handler" -t bugfix --tags redis,async

# List all prompts
pmt list

# Apply a prompt (interactive selector)
pmt apply

# Apply and delete (like git stash pop)
pmt pop

# Show prompt details
pmt show a7f

# Delete a prompt
pmt delete a7f
```

## Commands

### `pmt push <content>`

Save a new prompt to your local store.

**Options:**
- `-t, --type`: Type of prompt (bugfix, feature, refactor, general) - default: general
- `-g, --tags`: Comma-separated tags

**Examples:**
```bash
pmt push "Fix memory leak in async handler"
pmt push "Add OAuth login" -t feature --tags auth,api
pmt push "Refactor error handling" -t refactor
```

### `pmt list` (alias: `ls`)

List all saved prompts in a table format.

**Options:**
- `-t, --type`: Filter by type
- `-p, --project`: Filter by project

**Examples:**
```bash
pmt list
pmt list -t bugfix
pmt list -p my-api
pmt list -t feature -p my-api
```

### `pmt apply`

Interactively select a prompt and copy it to clipboard.

- Use ↑↓ arrow keys to navigate
- Press `/` to search
- Press Enter to select
- The prompt remains in storage after applying

**Example:**
```bash
pmt apply
```

### `pmt pop`

Interactively select a prompt, copy it to clipboard, and delete it from storage.

Similar to `git stash pop` - use this when you want to consume the prompt.

**Example:**
```bash
pmt pop
```

### `pmt show <id>`

Show detailed information about a specific prompt.

You can use the full ID or just a prefix.

**Examples:**
```bash
pmt show a7f3c2b
pmt show a7f
```

### `pmt delete <id>` (alias: `rm`)

Delete a specific prompt by its ID.

**Options:**
- `-f, --force`: Force deletion without confirmation

**Examples:**
```bash
pmt delete a7f3c2b
pmt delete a7f
pmt delete a7f -f  # Force delete without confirmation
```

## Storage

Prompts are stored in `~/.pmt/prompts.yaml`

You can back up your prompts by adding this directory to git:

```bash
cd ~/.pmt
git init
git add prompts.yaml
git commit -m "Backup prompts"
```

## Project Structure

```
pmt/
├── cmd/                  # Command implementations
│   ├── root.go          # Root command
│   ├── push.go          # Save prompts
│   ├── list.go          # List prompts
│   ├── apply.go         # Copy prompt
│   ├── pop.go           # Copy and delete
│   ├── show.go          # Show details
│   └── delete.go        # Delete prompt
├── internal/
│   ├── models/          # Data structures
│   │   └── prompt.go
│   ├── storage/         # Storage layer
│   │   └── store.go
│   ├── ui/              # Interactive UI
│   │   └── selector.go
│   └── utils/           # Utilities
│       ├── git.go       # Git detection
│       └── id.go        # ID generation
├── go.mod
├── go.sum
├── main.go
└── README.md
```

## Usage Example

```bash
# In your project directory
$ cd ~/projects/my-api

# Save a prompt
$ pmt push "Fix Redis connection leak in worker pool" -t bugfix --tags redis,async
✓ Saved prompt: a7f3c2b (bugfix) in project: my-api

# List all prompts
$ pmt list
ID        Type       Project    Content                                  Created
a7f3c2b   bugfix     my-api     Fix Redis connection leak in worker...   2025-01-13 15:30
9d4e1a8   feature    my-api     Add OAuth login with GitHub              2025-01-13 14:20

# Apply a prompt
$ pmt apply
? Select Prompt:
▸ a7f3c2b (bugfix) Fix Redis connection leak in worker pool
  9d4e1a8 (feature) Add OAuth login with GitHub

--------- Details ----------
ID:       a7f3c2b
Type:     bugfix
Project:  my-api
Created:  2025-01-13 15:30
Tags:     redis, async
Content:
Fix Redis connection leak in worker pool...

✓ Copied to clipboard: a7f3c2b
💡 Now paste (Ctrl+V) into Copilot!

# Show details
$ pmt show a7f
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
ID:        a7f3c2b
Type:      bugfix
Project:   my-api
Tags:      redis, async
Created:   2025-01-13 15:30:45
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Content:
Fix Redis connection leak in worker pool...

# Delete a prompt
$ pmt delete a7f
Delete prompt a7f3c2b? (Fix Redis connection leak in worker...)
Type 'yes' to confirm: yes
✓ Deleted prompt: a7f3c2b
```

## Requirements

- Go 1.21 or higher
- Git (for project detection)

## Dependencies

- [cobra](https://github.com/spf13/cobra) - CLI framework
- [promptui](https://github.com/manifoldco/promptui) - Interactive UI
- [clipboard](https://github.com/atotto/clipboard) - Clipboard operations
- [yaml.v3](https://gopkg.in/yaml.v3) - YAML parsing

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License - feel free to use this tool in your projects.

## Tips

1. Use descriptive prompts that capture context
2. Tag prompts for easy filtering
3. Use `pmt apply` to keep prompts, `pmt pop` to consume them
4. Back up your `~/.pmt/` directory
5. Use ID prefixes for quick access (e.g., `a7f` instead of `a7f3c2b`)

## Roadmap

Future enhancements:
- Search command with keyword matching
- Edit command for modifying prompts
- Statistics and usage tracking
- Import/export functionality
- Template support
- Remote sync capabilities