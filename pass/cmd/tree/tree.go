// Package tree provides tree view rendering for hierarchical data.
package tree

import (
	"strings"
)

// TreeNode represents a node in a hierarchical tree structure.
type TreeNode struct {
	Name     string
	IsDir    bool
	Children []*TreeNode
}

// NewTreeNode creates a new tree node with the given name and directory flag.
func NewTreeNode(name string, isDir bool) *TreeNode {
	return &TreeNode{
		Name:     name,
		IsDir:    isDir,
		Children: []*TreeNode{},
	}
}

// AddChild adds a child node and maintains alphabetical sort order.
func (n *TreeNode) AddChild(child *TreeNode) {
	n.Children = append(n.Children, child)
	// Sort children alphabetically using bubble sort (small datasets expected)
	for i := 0; i < len(n.Children)-1; i++ {
		for j := 0; j < len(n.Children)-i-1; j++ {
			if n.Children[j].Name > n.Children[j+1].Name {
				n.Children[j], n.Children[j+1] = n.Children[j+1], n.Children[j]
			}
		}
	}
}

// FindOrCreateChild finds a child by name or creates it if not found.
// Returns the child node (existing or newly created).
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

// Render renders the tree node and its children with box-drawing characters.
// The prefix parameter contains the indentation and vertical connectors for the current level.
func (n *TreeNode) Render(prefix string) string {
	var sb strings.Builder

	// Determine the connector: use └── for last child, ├── otherwise
	connector := "\u2514\u2500\u2500 " // └──
	if len(prefix) > 0 && !strings.HasSuffix(prefix, "    ") {
		connector = "\u251C\u2500\u2500 " // ├──
	}

	// Format the display name
	displayName := n.Name
	if n.IsDir {
		displayName += "/"
	}

	// Write the current node line
	sb.WriteString(prefix + connector + displayName + "\n")

	// Render children
	for i, child := range n.Children {
		isLast := i == len(n.Children)-1
		childPrefix := prefix
		if !isLast {
			childPrefix += "\u2502   " // │
		} else {
			childPrefix += "    "
		}
		sb.WriteString(child.Render(childPrefix))
	}

	return sb.String()
}

// BuildTreeFromPaths creates a tree structure from a list of file paths.
// Paths should NOT include the .gpg extension (they will be stripped if present).
// Returns a root node containing the tree structure.
func BuildTreeFromPaths(paths []string) *TreeNode {
	root := NewTreeNode("", false)

	for _, path := range paths {
		// Remove .gpg extension if present
		path = strings.TrimSuffix(path, ".gpg")

		// Split into path components
		components := strings.Split(path, "/")

		// Build tree structure from components
		current := root
		for i, component := range components {
			// Only the last component is a file, others are directories
			isDir := i < len(components)-1
			current = current.FindOrCreateChild(component, isDir)
		}
	}

	return root
}
