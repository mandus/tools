# TUI (Terminal User Interface)

The pass program uses Bubble Tea for an interactive TUI when fuzzy search is needed.

## When TUI is Used

The TUI automatically activates for:
- `pass` (no arguments) - Show mode
- `pass -c` (no arguments) - Clip mode  
- `pass rm` (no arguments) - Remove mode

## Keyboard Controls

| Key | Action |
|-----|--------|
| ↑ / ↓ | Navigate up/down |
| j / k | Navigate up/down |
| Enter | Select password |
| Esc | Exit/cancel |
| Ctrl+C | Exit/cancel |
| Ctrl+D | Exit/cancel |
| Ctrl+Q | Exit/cancel |
| Home | Jump to first item |
| End | Jump to last item |
| PageUp / PageDown | Scroll by page |
| Type | Filter passwords |
| Backspace | Delete character |

## Features

- Full terminal window utilization
- Real-time filtering as you type
- Match highlighting
- Clear help information
- Cross-platform compatibility

## Implementation

- Framework: [charmbracelet/bubbletea](https://github.com/charmbracelet/bubbletea)
- Components: [charmbracelet/bubbles](https://github.com/charmbracelet/bubbles)
- Styling: [charmbracelet/lipgloss](https://github.com/charmbracelet/lipgloss)
- Code: `cmd/tui/` package
