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
			wantCount: 4, // email/, gmail.com, social/, twitter.com
			contains:  []string{"gmail.com", "twitter.com", "email/", "social/"},
		},
		{
			name:     "Nested paths",
			passwords: []string{"email/gmail.com", "email/work/vpn.com"},
			wantCount: 4, // email/, gmail.com, work/, vpn.com
			contains:  []string{"gmail.com", "vpn.com", "email/", "work/"},
		},
		{
			name:     "Tree characters",
			passwords: []string{"a/1.com", "a/2.com", "b/3.com"},
			wantCount: 5, // a/, 1.com, 2.com, b/, 3.com
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

	// Initially should show all tree nodes (directories + password files)
	// email/gmail.com -> email/, gmail.com
	// email/work.com -> email/, work.com
	// social/twitter.com -> social/, twitter.com
	// Total: email/, gmail.com, work.com, social/, twitter.com = 5 items
	initialItems := model.list.Items()
	expectedCount := 5
	if len(initialItems) != expectedCount {
		t.Errorf("Expected %d initial items, got %d", expectedCount, len(initialItems))
	}

	// Set a query that should match gmail
	model.input.SetValue("gmail")
	model.filterList()

	// Should show email/ and gmail.com (2 items)
	filtered := model.list.Items()
	if len(filtered) != 2 {
		t.Errorf("Expected 2 filtered items for 'gmail' (email/ and gmail.com), got %d", len(filtered))
		return
	}

	// The filtered items should have gmail in their path
	for i, item := range filtered {
		path := item.FilterValue()
		if !strings.Contains(path, "gmail") && !strings.Contains(path, "email") {
			t.Errorf("Filtered item %d path %q doesn't contain 'gmail' or 'email'", i, path)
		}
	}
}

// TestTreeViewSelection tests that selection works with tree view
func TestTreeViewSelection(t *testing.T) {
	passwords := []string{"email/gmail.com", "social/twitter.com"}
	model := NewModel(passwords, FuzzyModeShow)

	// The tree structure will be:
	// email/ -> gmail.com
	// social/ -> twitter.com
	// So items are: email/, gmail.com, social/, twitter.com

	// Select the first item (email/)
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

	// The path should be one of the tree nodes (either a directory or a password)
	// Since we're showing all nodes in the tree, it could be a directory path
	validPaths := []string{"email", "email/gmail.com", "social", "social/twitter.com"}
	validPath := false
	for _, p := range validPaths {
		if path == p {
			validPath = true
			break
		}
	}
	if !validPath {
		t.Errorf("Selected path %q is not one of the valid tree node paths", path)
	}
}
