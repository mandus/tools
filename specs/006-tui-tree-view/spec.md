# Pass TUI Tree View Specification

## Overview

This document specifies the implementation of tree-style view in the interactive fuzzy finder TUI for all modes (show, clip, rm, edit), matching the visual hierarchy of the password store directory structure. The `pass find` command should remain with flat view as currently implemented.

## Status
- **Status**: Implemented ✅
- **Author**: @aasmundo
- **Created**: 2026-06-22
- **Last Updated**: 2026-06-22
- **Branch**: `feat/14-tui-tree-view`
- **Feature Number**: 006

## Background

Currently, the interactive fuzzy finder TUI displays passwords in a flat list format:
```
> email/gmail.com
  email/work.com
  social/twitter.com
```

This specification adds tree-based view to the TUI, similar to the `pass find` command output, to provide better visual hierarchy:
```
> email/
  ├── gmail.com
  └── work.com
  social/
  └── twitter.com
```

The `pass find <string>` command should continue to use flat view as currently implemented.

## Goals

- Display passwords in tree structure in the TUI fuzzy finder
- Maintain flat view for `pass find` command
- Support tree view in all TUI modes (show, clip, rm, edit)
- Maintain backward compatibility
- Follow existing tree rendering patterns from feature 004

## Non-Goals

- Changing the `pass find` command output (should remain flat)
- Adding collapsible/expandable tree nodes
- Adding color styling to tree view (can be future enhancement)
- Changing the fuzzy matching algorithm

## User Stories

### As a pass TUI user, I want to see passwords in a tree structure
So that I can better understand the directory hierarchy of my password store.

**Acceptance Criteria**:
- [x] TUI fuzzy finder displays passwords in tree structure
- [x] Tree structure accurately reflects the directory hierarchy
- [x] Tree uses box-drawing characters (├──, └──, │)
- [x] Directories are shown with trailing `/`
- [x] Files are shown without `.gpg` extension
- [x] Results are sorted alphabetically at each level
- [x] Works in all TUI modes (show, clip, rm, edit)
- [x] Selecting a directory with Enter automatically selects the first password within it

### As a pass CLI user, I want `pass find` to maintain flat output
So that my existing scripts and workflows continue to work.

**Acceptance Criteria**:
- [ ] `pass find <string>` continues to output flat list
- [ ] `pass find ls` outputs flat list (as before)
- [ ] No changes to existing CLI behavior

## Technical Design

### Architecture

The implementation will reuse the existing tree rendering logic from feature 004 (`cmd/tree/` package) and integrate it into the TUI fuzzy finder.

**Component Diagram:**
```
pass/
├── cmd/
│   ├── tui/
│   │   ├── models.go     # MODIFIED: Add tree view rendering
│   │   └── fuzzy.go      # MODIFIED: Pass tree structure to model
│   └── tree/             # EXISTING: Tree rendering package
│       ├── tree.go       # Tree node structure and rendering
│       └── tree_test.go  # Tree tests
└── specs/
    └── 006-tui-tree-view/
        ├── spec.md        # This document
        └── tasks.md       # Implementation tasks
```

### Data Flow

1. **Current Flow (Flat View):**
   ```
   CollectAllPasswords() → []string → TUI List Items
   ```

2. **New Flow (Tree View):**
   ```
   CollectAllPasswords() → []string → BuildTreeFromPaths() → TreeNode → RenderTreeForTUI() → TUI List Items
   ```

### Tree Node Structure

Reuse existing `TreeNode` from `cmd/tree/tree.go`:
```go
type TreeNode struct {
    Name     string
    IsDir    bool
    Children []*TreeNode
}
```

### TUI Integration

**Modified Model struct:**
```go
type Model struct {
    // Existing fields...
    allPasswords []string           // Keep for filtering
    treeRoot     *tree.TreeNode     // NEW: Tree structure for display
    flatView     bool               // NEW: Flag to toggle between views
}
```

**Tree Construction:**
- Build tree from `allPasswords` using `tree.BuildTreeFromPaths()`
- Store tree root in model
- Use tree structure for rendering list items

**Tree to List Conversion:**
Since Bubble Tea's `list` component expects a flat list of items, we need to flatten the tree for display while maintaining the visual hierarchy.

```go
// TreeItem represents a tree node as a list item
type TreeItem struct {
    path         string
    displayName  string  // e.g., "├── email/" or "└── gmail.com"
    indent       string   // Indentation prefix
    matchIndices []int    // For fuzzy match highlighting
}

// Implement list.Item interface
func (t TreeItem) Title() string       { return t.displayName }
func (t TreeItem) Description() string { return "" }
func (t TreeItem) FilterValue() string { return t.path }
```

### Rendering Strategy

**Option A: Flatten Tree to List Items (Recommended)**
- Convert tree structure to flat list of `TreeItem` with proper indentation
- Each tree node becomes a list item with formatted display name
- Maintains compatibility with existing list component
- Preserves fuzzy matching on full path

**Option B: Custom List Delegate with Tree Rendering**
- Keep existing item structure
- Modify delegate to render tree structure
- More complex but potentially more flexible

**Decision: Option A** - Simpler to implement and maintain compatibility.

### Tree to List Conversion Algorithm

```go
func flattenTreeToListItems(root *tree.TreeNode, prefix string) []list.Item {
    var items []list.Item
    
    for i, child := range root.Children {
        // Determine connector
        connector := "└── "
        if i < len(root.Children)-1 {
            connector = "├── "
        }
        
        // Build display name
        displayName := connector
        if child.IsDir {
            displayName += child.Name + "/"
        } else {
            displayName += child.Name
        }
        
        // Full path for filtering
        fullPath := buildFullPath(child)
        
        // Create tree item
        item := TreeItem{
            path:        fullPath,
            displayName: prefix + displayName,
            indent:      prefix,
        }
        items = append(items, item)
        
        // Recursively process children
        childPrefix := prefix
        if i < len(root.Children)-1 {
            childPrefix += "│   "
        } else {
            childPrefix += "    "
        }
        items = append(items, flattenTreeToListItems(child, childPrefix)...)
    }
    
    return items
}
```

### Filtering Considerations

- Fuzzy matching should still work on the full path (e.g., `email/gmail.com`)
- Display name includes tree characters but filtering uses clean path
- Match highlighting should work on the actual path, not the display name

### Mode Support

Tree view should be enabled for all TUI modes:
- `FuzzyModeShow` - Show password
- `FuzzyModeClip` - Copy to clipboard
- `FuzzyModeRm` - Remove password
- `FuzzyModeEdit` - Edit password

## Implementation Plan

### Phase 1: Test-Driven Development
1. [ ] Create tests for tree view rendering in TUI
2. [ ] Create tests for tree to list conversion
3. [ ] Create tests for filtering with tree view

### Phase 2: Core Implementation
1. [ ] Add tree construction to TUI model
2. [ ] Implement tree to list conversion
3. [ ] Update list delegate for tree rendering
4. [ ] Modify filtering to work with tree structure

### Phase 3: Integration
1. [ ] Integrate tree view into all TUI modes
2. [ ] Ensure backward compatibility
3. [ ] Add configuration option (if needed)

### Phase 4: Testing
1. [ ] Run all existing tests
2. [ ] Run new tree view tests
3. [ ] Manual testing of all modes
4. [ ] Verify `pass find` still uses flat view

## Testing Strategy

### Unit Tests

**New test file: `cmd/tui/tree_test.go`**

```go
package tui

import (
    "testing"
    "github.com/mandu/tools/pass/cmd/tree"
)

func TestFlattenTreeToListItems(t *testing.T) {
    // Build test tree
    root := tree.NewTreeNode("", false)
    email := tree.NewTreeNode("email", true)
    gmail := tree.NewTreeNode("gmail.com", false)
    work := tree.NewTreeNode("work.com", false)
    
    email.AddChild(gmail)
    email.AddChild(work)
    root.AddChild(email)
    
    // Flatten to list items
    items := flattenTreeToListItems(root, "")
    
    // Should have 3 items: email/, ├── gmail.com, └── work.com
    if len(items) != 3 {
        t.Errorf("Expected 3 items, got %d", len(items))
    }
    
    // Check first item
    if items[0].(TreeItem).displayName != "email/" {
        t.Errorf("Expected 'email/', got '%s'", items[0].(TreeItem).displayName)
    }
}

func TestTreeViewFiltering(t *testing.T) {
    // Test that filtering works with tree view
    // Match should find items by full path, not display name
}

func TestTreeViewSelection(t *testing.T) {
    // Test that selecting a tree item returns the correct path
}
```

### Integration Tests

**Modify existing `cmd/tui/tui_test.go`:**
- Add tests for tree view in different modes
- Verify flat view is no longer used

### Manual Testing

- Test tree view in show mode
- Test tree view in clip mode
- Test tree view in rm mode
- Test tree view in edit mode
- Test with various directory structures
- Test fuzzy matching with tree view
- Verify `pass find` still uses flat view

## Compatibility

### Backward Compatibility

- ✅ All existing CLI commands work unchanged
- ✅ `pass find` command maintains flat output
- ✅ All TUI keyboard shortcuts work unchanged
- ✅ Fuzzy matching algorithm unchanged

### Forward Compatibility

- Tree structure can be extended for future features (collapsible nodes, icons, etc.)
- Flat view can be added as an option if users prefer it

## Configuration

### Optional: Add Tree View Toggle

If users prefer flat view in TUI, we can add a configuration option:

```go
// In config package
var TreeViewEnabled = true  // Default to tree view

// Can be set via environment variable
// PASS_TREE_VIEW=false for flat view
```

**Decision**: Start with tree view as default, can add toggle later if needed.

## Error Handling

| Scenario | Behavior |
|----------|----------|
| Empty password store | Show empty list (current behavior) |
| Single password | Show as single item (no tree structure needed) |
| Tree construction error | Fall back to flat view with warning |
| Rendering error | Fall back to flat view |

## Open Questions

### OQ-001: Should we add a flat view option in TUI?
**Status**: OPEN - Propose: Not for initial implementation, can add later if users request it

### OQ-002: How to handle very long paths in tree view?
**Status**: OPEN - Propose: Let list component handle truncation as it does now

### OQ-003: Should directories be selectable in tree view?
**Status**: **RESOLVED** - Directories are displayed in the tree structure but when selected (Enter key), the TUI automatically selects the first password file within that directory. This allows users to see the full hierarchy while ensuring only actual password files are actionable.

## Success Criteria

- [ ] TUI displays passwords in tree structure
- [ ] Tree structure accurately reflects directory hierarchy
- [ ] All TUI modes support tree view
- [ ] `pass find` command maintains flat output
- [ ] All tests pass
- [ ] No breaking changes to existing functionality
- [ ] Documentation updated

## Appendix

### Example Tree Output in TUI

```
Select password (Enter to show, Esc to cancel)

> email/
  ├── gmail.com
  └── work.com
  social/
  ├── twitter.com
  └── github.com
  work/
  └── vpn/
      └── company.com

Search: 

↑/↓: Navigate | Enter: Select | Esc/Ctrl+C: Cancel | Ctrl+Q: Quit
```

### Comparison with Current Flat View

**Current (Flat):**
```
> email/gmail.com
  email/work.com
  social/twitter.com
  social/github.com
  work/vpn/company.com
```

**New (Tree):**
```
> email/
  ├── gmail.com
  └── work.com
  social/
  ├── twitter.com
  └── github.com
  work/
  └── vpn/
      └── company.com
```

## References

- [Feature 004: Pass Tree View for Find Command](../004-pass-tree-view/spec.md)
- [Bubble Tea List Component](https://github.com/charmbracelet/bubbles/tree/master/list)
- [Existing TUI Implementation](../pass-tui-spec.md)

---

## Implementation Summary

The tree view feature has been successfully implemented for the pass TUI fuzzy finder.

### Changes Made

1. **TUI Tree View**: Added tree-formatted display to the interactive fuzzy finder
2. **CLI Flat View**: Modified `pass find` command to use flat view by default
3. **Tree Formatting**: Reused existing `cmd/tree` package for tree rendering
4. **Backward Compatibility**: All existing functionality preserved
5. **Directory Selection**: When Enter is pressed on a directory, automatically selects the first password file within that directory
6. **Item Metadata**: Added `isDir` and `isPassword` fields to distinguish between directory and password nodes

### Files Modified

- `cmd/tui/models.go`: Added tree formatting support
- `cmd/find.go`: Changed to flat view by default
- `cmd/tui/tui_test.go`: Updated for new item types
- `cmd/tui/tree_view_test.go`: NEW - Tree view tests

### Testing

All tests pass, including new tree view tests:
- Tree formatting for various path structures
- Path preservation for filtering
- Filtering with tree view
- Selection with tree view

*Document Version: 1.1*
*Last Updated: 2026-06-22*
*Author: @aasmundo*
*Status: Implemented*
*Feature Number: 006*