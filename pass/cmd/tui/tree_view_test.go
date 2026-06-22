package tui

import (
	"strings"
	"testing"
)

// TestCreateTreeFormattedItems tests the tree formatting of password paths
func TestCreateTreeFormattedItems(t *testing.T) {
	tests := []struct {
		name     string
		passwords []string
		wantCount int
		contains  []string // Strings that should be in the display names
	}{
		{
			name:     "Single level paths",
			passwords: []string{"email/gmail.com", "social/twitter.com"},
			wantCount: 2,
			contains:  []string{"gmail.com", "twitter.com"},
		},
		{
			name:     "Nested paths",
			passwords: []string{"email/gmail.com", "email/work/vpn.com"},
			wantCount: 2,
			contains:  []string{"gmail.com", "vpn.com"},
		},
		{
			name:     "Tree characters",
			passwords: []string{"a/1.com", "a/2.com", "b/3.com"},
			wantCount: 3,
			contains:  []string{"├──", "└──"}, // Should contain tree characters
		},
		{
			name:     "Empty list",
			passwords: []string{},
			wantCount: 0,
			contains:  []string{},
		},
		{
			name:     "Single password",
			passwords: []string{"single.com"},
			wantCount: 1,
			contains:  []string{"single.com"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			items := CreateTreeFormattedItems(tt.passwords)

			// Check count
			if len(items) != tt.wantCount {
				t.Errorf("Expected %d items, got %d", tt.wantCount, len(items))
				return
			}

			// Check that display names contain expected strings
			for _, expected := range tt.contains {
				found := false
				for _, listItem := range items {
					// Try to get the title using type assertion
					var title string
					if i, ok := listItem.(item); ok {
						title = i.path // For item type, path is the title
					} else if ti, ok := listItem.(treeFormattedItem); ok {
						title = ti.displayName // For treeFormattedItem, use displayName
					} else {
						continue
					}
					if strings.Contains(title, expected) {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected to find %q in display names", expected)
				}
			}
		})
	}
}

// TestTreeViewPreservesPaths tests that the full paths are preserved for filtering
func TestTreeViewPreservesPaths(t *testing.T) {
	passwords := []string{"email/gmail.com", "social/twitter.com"}
	items := CreateTreeFormattedItems(passwords)

	// Check that we can retrieve the original paths
	foundPaths := make(map[string]bool)
	for _, item := range items {
		path := item.FilterValue()
		foundPaths[path] = true
	}

	// All original paths should be present
	for _, path := range passwords {
		if !foundPaths[path] {
			t.Errorf("Original path %q not found in tree items", path)
		}
	}
}

// TestTreeViewFiltering tests that filtering works with tree view
func TestTreeViewFiltering(t *testing.T) {
	passwords := []string{"email/gmail.com", "email/work.com", "social/twitter.com"}
	
	// Create model with tree view
	model := NewModel(passwords, FuzzyModeShow)

	// Initially should show all items
	initialItems := model.list.Items()
	if len(initialItems) != len(passwords) {
		t.Errorf("Expected %d initial items, got %d", len(passwords), len(initialItems))
	}

	// Set a query that should match gmail
	model.input.SetValue("gmail")
	model.filterList()

	// Should only show gmail item
	filtered := model.list.Items()
	if len(filtered) != 1 {
		t.Errorf("Expected 1 filtered item for 'gmail', got %d", len(filtered))
		return
	}

	// The filtered item should have gmail in its path
	path := filtered[0].FilterValue()
	if !strings.Contains(path, "gmail") {
		t.Errorf("Filtered item path %q doesn't contain 'gmail'", path)
	}
}

// TestTreeViewSelection tests that selection works with tree view
func TestTreeViewSelection(t *testing.T) {
	passwords := []string{"email/gmail.com", "social/twitter.com"}
	model := NewModel(passwords, FuzzyModeShow)

	// Select the first item
	model.list.Select(0)

	// Get selected item
	selectedItem := model.list.SelectedItem()
	if selectedItem == nil {
		t.Fatal("No item selected")
	}

	// The selected item should have a valid path
	path := selectedItem.FilterValue()
	if path == "" {
		t.Error("Selected item has empty path")
	}

	// The path should be one of the original passwords
	validPath := false
	for _, p := range passwords {
		if path == p {
			validPath = true
			break
		}
	}
	if !validPath {
		t.Errorf("Selected path %q is not one of the original passwords", path)
	}
}
