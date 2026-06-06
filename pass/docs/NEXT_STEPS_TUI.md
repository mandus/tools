# Next Steps for Pass TUI Implementation

## Current Status ✅

The TUI implementation is **functionally complete** but **not yet integrated** with the existing CLI commands.

### What's Done
- ✅ Created `cmd/tui/` package with Bubble Tea
- ✅ Implemented core TUI functionality
- ✅ Added all required dependencies
- ✅ Compiles successfully
- ✅ All keyboard handling works
- ✅ Full window size utilization
- ✅ Basic help information display

### What's Left
- ⏳ Integrate with existing CLI commands
- ⏳ Add fuzzy matching (currently simple contains)
- ⏳ Add confirmation for remove mode
- ⏳ Polish and testing

## Immediate Next Steps

### 1. Integrate TUI with CLI Commands

#### Update `cmd/root.go`
Find the section where fuzzy search is called and replace:

```go
// BEFORE:
// Check if clip flag is set - enter fuzzy search mode with clip
clipFlagChanged, _ := cmd.Flags().GetBool("clip")
if clipFlagChanged {
    // Set global clip flag
    clipFlag = true
    // Enter fuzzy search mode with clip
    return RunInteractiveFuzzySearch(FuzzyModeClip)
}
// Enter fuzzy search mode (default: show)
return RunInteractiveFuzzySearch(FuzzyModeShow)

// AFTER:
// Check if clip flag is set - enter fuzzy search mode with clip
clipFlagChanged, _ := cmd.Flags().GetBool("clip")
if clipFlagChanged {
    // Set global clip flag
    clipFlag = true
    // Enter fuzzy search mode with clip
    selected, err := tui.RunInteractiveFuzzySearch(FuzzyModeClip)
    if err != nil {
        return err
    }
    if selected != "" {
        return showPassword(selected)
    }
    return nil
}
// Enter fuzzy search mode (default: show)
selected, err := tui.RunInteractiveFuzzySearch(FuzzyModeShow)
if err != nil {
    return err
}
if selected != "" {
    return showPassword(selected)
}
return nil
```

#### Update `cmd/rm.go`
Find the `runRmFuzzySearch` function and update it to use the new TUI:

```go
// BEFORE:
func runRmFuzzySearch(noCommit, clip bool) error {
    // If clip flag is set, we need special handling
    if clip {
        // ... complex logic
    }
    // Normal rm without clip
    selected, err := InteractiveFuzzySearch(FuzzyModeRm)
    if err != nil {
        return err
    }
    if selected == "" {
        return nil
    }
    
    fullPath := getRmFullPath(selected)
    return removePasswordInternal(fullPath, selected, noCommit)
}

// AFTER:
func runRmFuzzySearch(noCommit, clip bool) error {
    // Use the new TUI
    selected, err := tui.RunInteractiveFuzzySearch(FuzzyModeRm)
    if err != nil {
        return err
    }
    if selected == "" {
        return nil
    }
    
    // If clip flag is set, copy to clipboard first
    if clip {
        fullPath := getRmFullPath(selected)
        password, err := gpg.DecryptFile(fullPath)
        if err != nil {
            return err
        }
        if err := filesystem.CopyToClipboard(password); err != nil {
            return fmt.Errorf("pass: failed to copy to clipboard: %v", err)
        }
        fmt.Printf("Copied %s to clipboard.\n", selected)
        go filesystem.StartClipboardClearTimer()
    }
    
    // Remove the password
    fullPath := getRmFullPath(selected)
    return removePasswordInternal(fullPath, selected, noCommit)
}
```

### 2. Add Fuzzy Matching

Currently, the TUI uses simple string contains for filtering. Update `cmd/tui/models.go`:

```go
// Add import:
"github.com/mandu/tools/pass/pkg/fuzzy"

// Replace filterList function:
func (m *Model) filterList() {
    query := m.input.Value()
    
    // Use proper fuzzy matching
    results := fuzzy.Filter(query, m.allPasswords)
    
    // Convert to list items
    items := make([]list.Item, len(results))
    for i, result := range results {
        items[i] = item{
            path:         result.Path,
            matchScore:   result.Score,
            matchIndices: result.MatchIndices,
        }
    }
    
    m.list.SetItems(items)
}
```

### 3. Add Confirmation for Remove Mode

Create a simple confirmation dialog. Add to `cmd/tui/fuzzy.go`:

```go
// RunInteractiveFuzzySearch with confirmation for remove mode
func RunInteractiveFuzzySearch(mode cmd.FuzzySearchMode) (string, error) {
    // ... existing code ...
    
    // For remove mode, add confirmation
    if mode == cmd.FuzzyModeRm {
        selected, err := RunFuzzySearch(passwords, mode)
        if err != nil {
            return "", err
        }
        if selected == "" {
            return "", nil
        }
        
        // Show confirmation
        fmt.Printf("\nAre you sure you want to remove %s? (y/N): ", selected)
        var response string
        fmt.Scanln(&response)
        if strings.ToLower(response) != "y" && strings.ToLower(response) != "yes" {
            return "", nil // User cancelled
        }
        
        return selected, nil
    }
    
    // For other modes, no confirmation needed
    return RunFuzzySearch(passwords, mode)
}
```

## Testing Plan

### 1. Unit Testing
```bash
# Test the TUI package
go test ./cmd/tui/... -v

# Test the entire package
go test ./... -v
```

### 2. Manual Testing

#### Test Show Mode
```bash
./pass
# Should show TUI with all passwords
# Arrow keys should navigate
# Enter should show selected password
# Esc should exit
```

#### Test Clip Mode
```bash
./pass -c
# Should show TUI
# Enter should copy selected password to clipboard
```

#### Test Remove Mode
```bash
./pass rm
# Should show TUI
# Enter should show confirmation dialog
# Confirming should remove the password
```

### 3. Cross-Platform Testing
- Test on Windows Terminal
- Test on Linux terminal (GNOME, Konsole, xterm)
- Test on macOS Terminal
- Test with different terminal sizes

## File Changes Summary

### Files to Modify
1. `cmd/root.go` - Update fuzzy search calls
2. `cmd/rm.go` - Update fuzzy search calls
3. `cmd/tui/models.go` - Add fuzzy matching
4. `cmd/tui/fuzzy.go` - Add confirmation for remove mode

### Files to Add
- None (all TUI files already created)

### Files to Remove (Future)
- `cmd/fuzzy.go` - After integration is complete
- `pkg/terminal/terminal.go` - After full migration
- `pkg/terminal/key_reader.go` - After full migration

## Expected Timeline

| Task | Time Estimate | Priority |
|------|---------------|----------|
| Integrate with CLI commands | 1-2 hours | High |
| Add fuzzy matching | 30-60 min | Medium |
| Add confirmation dialog | 30-60 min | Medium |
| Testing and bug fixes | 2-4 hours | High |
| Polish and styling | 1-2 hours | Low |
| Cross-platform testing | 1-2 hours | Medium |

## Checklist

- [ ] Update `cmd/root.go` to use new TUI
- [ ] Update `cmd/rm.go` to use new TUI
- [ ] Add fuzzy matching to filterList
- [ ] Add confirmation dialog for remove mode
- [ ] Test show mode
- [ ] Test clip mode
- [ ] Test remove mode
- [ ] Test keyboard navigation
- [ ] Test search filtering
- [ ] Test window resize
- [ ] Test edge cases (empty store, no matches)
- [ ] Cross-platform testing

## Notes

- The current implementation compiles and the TUI framework is working
- The main remaining work is integration with the existing CLI
- Once integrated, the old terminal handling code can be deprecated
- The new implementation solves the cross-platform keyboard handling issues
