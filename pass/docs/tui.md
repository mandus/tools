# TUI (Terminal User Interface)

The pass program uses Bubble Tea for an interactive TUI when fuzzy search is needed.

## When TUI is Used

The TUI automatically activates for:
- `pass` (no arguments) - Show mode
- `pass -c` (no arguments) - Clip mode  
- `pass rm` (no arguments) - Remove mode
- `pass edit` (no arguments) - Edit mode

## Keyboard Controls

| Key | Action |
|-----|--------|
| ↑ / ↓ | Navigate list up/down |
| Enter | Select password (action depends on mode: show, copy, remove, or edit) |
| j / k | Navigate list up/down |
| ← / → | Move cursor in search input |
| Home / End | Move cursor to start/end of search input |
| Enter | Select password |
| Esc | Exit/cancel |
| Ctrl+C | Exit/cancel |
| Ctrl+D | Exit/cancel |
| Ctrl+Q | Exit/cancel |
| Home | Jump to first item |
| End | Jump to last item |
| PageUp / PageDown | Scroll by page |
| Tab | Cycle through results |
| Type | Filter passwords (fuzzy matching) |
| Backspace | Delete character |
| Ctrl+A | Move to start of query |
| Ctrl+E | Move to end of query |
| Ctrl+K | Clear from cursor to end |
| Ctrl+L | Clear entire query |
| Ctrl+W | Delete word before cursor |

## Modes

The TUI operates in different modes depending on how it was invoked:

| Mode | Invocation | Action on Enter |
|------|------------|----------------|
| Show | `pass` | Show the selected password |
| Clip | `pass -c` | Copy selected password to clipboard |
| Remove | `pass rm` | Delete the selected password |
| Edit | `pass edit` | Edit the selected password in $EDITOR |

## Fuzzy Matching

The TUI uses **subsequence matching** - characters must appear in order but not necessarily consecutively:

- Typing `"g"` → shows all items containing "g" anywhere
- Typing `"ga"` → shows items where "g" appears anywhere, followed by "a" anywhere after it
- Typing `"gmail"` → shows items like `email/gmail.com/user` where characters appear in order

Matching characters are **highlighted** in the results for better visibility.

## Features

- Full terminal window utilization
- Real-time fuzzy filtering as you type
- Match highlighting with fuzzy match indices
- Clear help information
- Cross-platform compatibility
- Proper arrow key navigation
- Tab key to cycle through results

## Implementation

- Framework: [charmbracelet/bubbletea](https://github.com/charmbracelet/bubbletea)
- Components: [charmbracelet/bubbles](https://github.com/charmbracelet/bubbles)
- Styling: [charmbracelet/lipgloss](https://github.com/charmbracelet/lipgloss)
- Code: `cmd/tui/` package
- Fuzzy matching: Uses the existing `pkg/fuzzy` package for subsequence matching
