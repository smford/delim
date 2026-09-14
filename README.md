# delim

A lightweight, high-performance terminal utility that prints visual separator lines across terminal screens.

[![CI](https://github.com/smford/delim/actions/workflows/ci.yaml/badge.svg)](https://github.com/smford/delim/actions/workflows/ci.yaml)
[![Release](https://github.com/smford/delim/actions/workflows/release.yaml/badge.svg)](https://github.com/smford/delim/actions/workflows/release.yaml)

## Features

- **Blazingly fast**: Built in Go with zero extra allocations on single-character lines and a single output write.
- **Dynamic terminal detection**: Automatically detects width from stdout, stderr, stdin, or `$COLUMNS`, with graceful fallback for pipes and non-interactive environments.
- **Configurable**: Configurable via flags, environment variables, or YAML config file.
- **Cross-platform**: Binaries available for macOS (Apple Silicon & Intel), Linux (AMD64 & ARM64), and Windows.

## Installation

### Homebrew

```bash
brew install smford/tap/delim
```

### Binary Download

Download pre-compiled binaries for macOS, Linux, and Windows from the [GitHub Releases](https://github.com/smford/delim/releases) page.

### Go Install

```bash
go install github.com/smford/delim@latest
```

### From Source

```bash
git clone https://github.com/smford/delim.git
cd delim
go build -o delim .
```

## Usage

Simply run:

```bash
delim
```

### Command Line Options

```text
Usage: delim [options]

Options:
  -c, --char string        Default line character or pattern (default "=")
      --config string      Configuration file: /path/to/file.yaml (default "~/.delim")
  -w, --width int          Line width (default: terminal width)
  -n, --newline            Print trailing newline (default false)
      --displayconfig      Display configuration
  -v, --version            Display version
  -h, --help               Display help
```

### Examples

```bash
# Print default '=' line matching terminal width
delim

# Print line of hyphens with trailing newline
delim -c "-" -n

# Print line of custom width (e.g. 40 characters)
delim -w 40 -n

# Use multi-character pattern
delim -c "=-" -w 30 -n
```

## Configuration

Precedence order: **Command Line Flags > Environment Variables > Config File > Defaults**.

### 1. Command Line Flags

```bash
delim --char "-"
```

### 2. Environment Variables

All environment variables use the `DELIM_` prefix:

```bash
export DELIM_CHAR="-"
export DELIM_WIDTH=80
export DELIM_NEWLINE=true
```

### 3. Configuration File

Configuration is optional. By default, `delim` checks for `~/.delim` (in YAML format):

```yaml
char: "-"
newline: false
width: 0
```

To display active settings:

```bash
delim --displayconfig
```

