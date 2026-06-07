# Pass Tree View for Find Command Specification

## Overview

This document specifies the implementation of a tree-style view for the `pass find` command, matching the visual hierarchy of the original Unix pass tool. Currently, `pass find` displays results in a flat list. This feature will display results in a tree structure that reflects the directory hierarchy of the password store.

### Current Behavior
```bash
$ pass find gmail
email/gmail.com
social/gmail-backup.com
work/gmail-work.com
```

### Required Behavior
```bash
$ pass find gmail
├── email/
│   └── gmail.com
├── social/
│   └── gmail-backup.com
└── work/
    └── gmail-work.com
```

---

## User Requirements

### Must Have
- [ ] `pass find <string>` displays results in a tree structure
- [ ] Tree structure accurately reflects the directory hierarchy from `~/.password-store/`
- [ ] Tree uses box-drawing characters (├──, └──, │) for visual hierarchy
- [ ] Directories are shown with trailing `/` 
- [ ] Files are shown without `.gpg` extension
- [ ] Results are sorted alphabetically at each level
- [ ] Indentation correctly shows nesting depth

### Should Have
- [ ] Support `--flat` flag to revert to original flat list output
- [ ] Support `--no-tree` flag as alias for `--flat`
- [ ] Color output for different elements (directories vs files)
- [ ] Consistent with existing `pass ls` command behavior

### Nice to Have (Future)
- [ ] Collapsible tree view in TUI mode
- [ ] Tree view for `pass ls` command
- [ ] Customizable tree characters

---

## Detailed Specifications

### 1. Tree Structure Rendering

**Tree Characters:**
```
├── : Last item in a group (has siblings below)
└── : Last item in a group (no siblings below)
│   : Vertical connector (parent has siblings below)
    : Indentation (4 spaces per level)
```

**Example with deep nesting:**
```
├── dev/
│   ├── hafslund/
│   │   └── mistral-vibe-key
│   └── mistral.ai/
│       ├── api-access-alternate-key
│       ├── asmund.odegard@hafslund.no
│       └── for-pi-api-key
└── nucmman/
    └── mistral-vibe-key
```

**Implementation Requirements:**
- Build a tree data structure from matched paths
- Sort children alphabetically at each level
- Determine correct prefix characters (├──, └──) based on position
- Indent each level with 4 spaces
- Directories must end with `/`
- Files must not include `.gpg` extension

### 2. Path Parsing and Tree Construction

**Input:** List of matched paths from existing find logic
```go
[]string{
    "dev/hafslund/mistral-vibe-key",
    "dev/mistral.ai/api-access-alternate-key",
    "dev/mistral.ai/asmund.odegard@hafslund.no",
    "dev/mistral.ai/for-pi-api-key",
    "nucmman/mistral-vibe-key",
}
```

**Tree Construction Algorithm:**
1. Parse each path into components (split by `/`)
2. Build a tree where each node represents a directory or file
3. Mark nodes as directory (if they have children or path continues) or file (leaf node)
4. Sort children alphabetically at each level

**Node Structure:**
```go
type TreeNode struct {
    Name     string    // e.g., "dev", "hafslund", "mistral-vibe-key"
    IsDir    bool      // true if this node is a directory
    Children []*TreeNode  // child nodes, sorted alphabetically
}
```

### 3. Tree Rendering Algorithm

**Render Function:**
```go
func renderTree(node *TreeNode, prefixes []string) string
```

**Parameters:**
- `node`: Current tree node to render
- `prefixes`: Array of prefix strings for each level (built during recursion)

**Logic:**
1. Determine current prefix based on position in parent's children
   - If last child: `"└── "`
   - Otherwise: `"├── "`
2. Build display name:
   - If IsDir: `name + "/"`
   - Otherwise: `name`
3. For each child:
   - Determine if child is last in its group
   - Build child prefix: parent prefix + ("│   " if parent has more siblings, "    " otherwise)
   - Recursively render child

**Example:**
```
root
└── dev/
    ├── hafslund/
    │   └── mistral-vibe-key
    └── mistral.ai/
        ├── api-access-alternate-key
        ├── asmund.odegard@hafslund.no
        └── for-pi-api-key
```

### 4. Command Integration

**Modified Command:**
```bash
pass find [OPTIONS] <string>
```

**New Flags:**
- `--flat, -f`: Display results as flat list (original behavior)
- `--no-tree`: Alias for `--flat`

**Flag Behavior:**
- Default (no flags): Tree view
- `--flat` or `--no-tree`: Flat list view (original behavior)

**Modified findCmd:**
```go
var findCmd = &cobra.Command{
    Use:   "find [string]",
    Short: "Search for passwords",
    Long:  `Search for passwords containing the given string anywhere in their path.`,
    Args:  cobra.ExactArgs(1),
    RunE: func(cmd *cobra.Command, args []string) error {
        flat, _ := cmd.Flags().GetBool("flat")
        return findPasswords(args[0], flat)
    },
}

func addFindCmd() {
    findCmd.Flags().BoolVarP(&flatFlag, "flat", "f", false, "Output flat list instead of tree")
    findCmd.Flags().BoolVar(&noTreeFlag, "no-tree", false, "Output flat list instead of tree")
    rootCmd.AddCommand(findCmd)
}
```

### 5. Tree View Styling (Optional Enhancement)

**Using Lip Gloss for styling:**
```go
var (
    dirStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#5DADE2"))
    fileStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF"))
    treeStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#7D3C98"))
)
```

**Styled Output:**
- Directories: Blue (#5DADE2)
- Files: White (#FFFFFF)  
- Tree characters: Purple (#7D3C98)

---

## Implementation Details

### 5.1 File Structure

```
pass/
├── cmd/
│   ├── find.go          # MODIFIED: Add tree view support
│   ├── find_test.go     # MODIFIED: Add tree view tests
│   └── tree/            # NEW: Tree rendering package
│       ├── tree.go      # Tree node structure and rendering
│       └── tree_test.go # Tree tests
└── specs/
    └── pass-tree-view/  # NEW: This spec
        ├── spec.md
        └── tasks.md
```

### 5.2 Tree Package Implementation

**tree.go:**
```go
package tree

import (
    "fmt"
    "path/filepath"
    "sort"
    "strings"
)

// TreeNode represents a node in the password store tree
type TreeNode struct {
    Name     string
    IsDir    bool
    Children []*TreeNode
}

// NewTreeNode creates a new tree node
func NewTreeNode(name string, isDir bool) *TreeNode {
    return &TreeNode{
        Name:     name,
        IsDir:    isDir,
        Children: []*TreeNode{},
    }
}

// AddChild adds a child node and maintains sorted order
func (n *TreeNode) AddChild(child *TreeNode) {
    n.Children = append(n.Children, child)
    sort.Slice(n.Children, func(i, j int) bool {
        return n.Children[i].Name < n.Children[j].Name
    })
}

// FindOrCreateChild finds a child by name or creates it if not found
func (n *TreeNode) FindOrCreateChild(name string, isDir bool) *TreeNode {
    for _, child := range n.Children {
        if child.Name == name {
            return child
        }
    }
    child := NewTreeNode(name, isDir)
    n.AddChild(child)
    return child
}

// Render renders the tree with box-drawing characters
func (n *TreeNode) Render(prefix string) string {
    var sb strings.Builder
    
    // Determine the connector
    connector := "└── "
    if len(prefix) > 0 && !strings.HasSuffix(prefix, "    ") {
        connector = "├── "
    }
    
    // Build display name
    displayName := n.Name
    if n.IsDir {
        displayName += "/"
    }
    
    sb.WriteString(prefix + connector + displayName + "\n")
    
    // Render children
    for i, child := range n.Children {
        isLast := i == len(n.Children)-1
        childPrefix := prefix
        if !isLast {
            childPrefix += "│   "
        } else {
            childPrefix += "    "
        }
        sb.WriteString(child.Render(childPrefix))
    }
    
    return sb.String()
}

// BuildTreeFromPaths builds a tree from a list of password paths
func BuildTreeFromPaths(paths []string) *TreeNode {
    root := NewTreeNode("", false)
    
    for _, path := range paths {
        // Remove .gpg extension if present
        path = strings.TrimSuffix(path, ".gpg")
        
        // Split into components
        components := strings.Split(path, "/")
        
        // Build tree structure
        current := root
        for i, component := range components {
            isDir := i < len(components)-1
            current = current.FindOrCreateChild(component, isDir)
        }
    }
    
    return root
}
```

### 5.3 Modified find.go

```go
package cmd

import (
    "fmt"
    "os"
    "path/filepath"
    "strings"

    "github.com/mandu/tools/pass/pkg/filesystem"
    "github.com/mandu/tools/pass/cmd/tree"
    "github.com/spf13/cobra"
)

// findCmd represents the find command
var findCmd = &cobra.Command{
    Use:   "find [string]",
    Short: "Search for passwords",
    Long:  `Search for passwords containing the given string anywhere in their path.`,
    Args:  cobra.ExactArgs(1),
    RunE: func(cmd *cobra.Command, args []string) error {
        flat, _ := cmd.Flags().GetBool("flat")
        noTree, _ := cmd.Flags().GetBool("no-tree")
        return findPasswords(args[0], flat || noTree)
    },
}

// Flags for find command
var (
    flatFlag   bool
    noTreeFlag bool
)

func addFindCmd() {
    findCmd.Flags().BoolVarP(&flatFlag, "flat", "f", false, "Output flat list instead of tree")
    findCmd.Flags().BoolVar(&noTreeFlag, "no-tree", false, "Output flat list instead of tree")
    rootCmd.AddCommand(findCmd)
}

// findPasswords searches for passwords containing the search string
func findPasswords(searchString string, flat bool) error {
    // ... existing search logic ...
    
    // Get matched paths
    var results []string
    // ... (existing filepath.Walk logic to populate results) ...
    
    if flat {
        // Original flat output
        for _, result := range results {
            fmt.Println(result)
        }
    } else {
        // New tree view
        root := tree.BuildTreeFromPaths(results)
        // Skip the root node (empty string) in rendering
        var output strings.Builder
        for i, child := range root.Children {
            prefix := ""
            if i < len(root.Children)-1 {
                prefix = "│   "
            } else {
                prefix = "    "
            }
            // Remove the leading connector from first level
            treeOutput := child.Render("")
            lines := strings.Split(treeOutput, "\n")
            for j, line := range lines {
                if j == 0 {
                    // First line already has the connector
                    output.WriteString(line + "\n")
                } else if strings.HasPrefix(line, "    ") {
                    // Indented lines - keep as is
                    output.WriteString(line + "\n")
                } else {
                    output.WriteString(line + "\n")
                }
            }
        }
        fmt.Print(output.String())
    }
    
    return nil
}
```

---

## Error Handling

| Scenario | Behavior |
|----------|----------|
| No matches found | Print nothing (current behavior) |
| Invalid search string | Error message (current behavior) |
| Store doesn't exist | Print nothing (current behavior) |
| Tree rendering error | Fall back to flat output with warning |

---

## Testing Strategy

### 6.1 Unit Tests

**tree/tree_test.go:**
```go
package tree

import (
    "testing"
)

func TestNewTreeNode(t *testing.T) {
    node := NewTreeNode("test", true)
    if node.Name != "test" {
        t.Errorf("Expected name 'test', got '%s'", node.Name)
    }
    if !node.IsDir {
        t.Error("Expected IsDir to be true")
    }
}

func TestAddChild(t *testing.T) {
    parent := NewTreeNode("parent", true)
    child1 := NewTreeNode("child1", false)
    child2 := NewTreeNode("child2", false)
    
    parent.AddChild(child2)
    parent.AddChild(child1) // Added second, should be sorted first
    
    if len(parent.Children) != 2 {
        t.Fatalf("Expected 2 children, got %d", len(parent.Children))
    }
    if parent.Children[0].Name != "child1" {
        t.Errorf("Expected first child to be 'child1', got '%s'", parent.Children[0].Name)
    }
}

func TestBuildTreeFromPaths(t *testing.T) {
    paths := []string{
        "email/gmail.com",
        "email/work.com",
        "social/twitter.com",
    }
    
    root := BuildTreeFromPaths(paths)
    
    if len(root.Children) != 2 {
        t.Fatalf("Expected 2 top-level children, got %d", len(root.Children))
    }
    
    // Check email directory
    emailDir := root.Children[0]
    if emailDir.Name != "email" {
        t.Errorf("Expected first child to be 'email', got '%s'", emailDir.Name)
    }
    if !emailDir.IsDir {
        t.Error("Expected email to be a directory")
    }
    if len(emailDir.Children) != 2 {
        t.Fatalf("Expected email to have 2 children, got %d", len(emailDir.Children))
    }
}

func TestRender(t *testing.T) {
    root := NewTreeNode("", false)
    email := NewTreeNode("email", true)
    gmail := NewTreeNode("gmail.com", false)
    
    email.AddChild(gmail)
    root.AddChild(email)
    
    output := email.Render("")
    expected := "email/\n    └── gmail.com\n"
    
    if output != expected {
        t.Errorf("Expected:\n%s\nGot:\n%s", expected, output)
    }
}
```

**cmd/find_test.go:**
- Test tree view output format
- Test flat view still works
- Test flags work correctly
- Test edge cases (empty results, single result, deep nesting)

### 6.2 Integration Tests

**tests/tree_view_test.go:**
```go
package tests

import (
    "testing"
    "os"
    "os/exec"
)

func TestFindTreeView(t *testing.T) {
    // Setup test store
    // Insert test passwords with nested structure
    // Run pass find
    // Verify tree output
}

func TestFindFlatView(t *testing.T) {
    // Setup test store
    // Run pass find --flat
    // Verify flat output
}
```

### 6.3 Manual Testing

- Test with various directory structures
- Test with single-level paths
- Test with deeply nested paths
- Test with special characters in paths
- Test with empty results
- Test flag combinations

---

## Compatibility

### Unix pass Compatibility

| Feature | Unix pass | This implementation |
|---------|-----------|-------------------|
| `pass find <string>` | Flat list | Tree view (default) |
| Tree structure | No | Yes |
| `--flat` flag | No | Yes |

### Cross-Platform Compatibility

- Tree characters (├──, └──, │) are standard ASCII/Unicode and work in all modern terminals
- No platform-specific code required
- Works with Windows Terminal, Linux terminals, macOS Terminal

---

## Environment Variables

No new environment variables required. Existing variables still apply:
- `PASSWORD_STORE_DIR`: Password store location

---

## Open Questions

### OQ-001: Should tree view be the default?
**Status**: DECIDED - Yes, tree view is more informative and matches user expectation from Unix pass screenshot

### OQ-002: Should we add color by default?
**Status**: OPEN - Propose: No for initial implementation, can be added as enhancement

### OQ-003: How to handle very deep nesting (>10 levels)?
**Status**: OPEN - Propose: No special handling, let terminal handle wrapping

---

*Feature Number: 004*

---

## References

- [Unix pass find command](https://git.zx2c4.com/password-store/tree/src/password-store.sh)
- [Bubble Tea tree rendering examples](https://github.com/charmbracelet/bubbletea/tree/master/examples)
- [Lip Gloss styling](https://github.com/charmbracelet/lipgloss)

---

*Document Version: 1.0*  
*Last Updated: 2026-06-07*  
*Author: Mandu*  
*Status: Approved for Implementation*
*Feature Number: 004*
