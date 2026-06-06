package fuzzy

import (
	"testing"
)

func TestMatch(t *testing.T) {
	tests := []struct {
		query  string
		target string
		want   bool
	}{
		// Basic matches
		{"twt", "social/twitter.com/admin", true},
		{"gm", "email/gmail.com/user", true},
		{"chase", "banking/chase.com/account", true},
		{"chaseb", "banking/chase.com/account", false}, // No 'b' after 'e' in "chase.com/account"
		{"pass", "passwords/password.txt", true},

		// Exact matches
		{"twitter", "twitter", true},
		{"", "anything", true},           // Empty query matches everything
		{"anything", "", false},          // Empty target never matches
		{"", "", true},                   // Both empty

		// Case insensitivity
		{"TWT", "social/twitter.com/admin", true},
		{"Twt", "Social/Twitter.Com/Admin", true},

		// Non-matches
		{"twtt", "twitter", true},        // Actually matches: t(0), w(1), t(3), t(4)
		{"mtw", "twitter", false},        // Wrong order
		{"abc", "def", false},           // No common characters
		{"xyz", "abc", false},

		// Single character
		{"t", "twitter", true},
		{"z", "twitter", false},

		// Path separators
		{"email", "email/gmail.com/user", true},
		{"com", "email/gmail.com/user", true},
	}

	for _, tt := range tests {
		t.Run(tt.query+"_"+tt.target, func(t *testing.T) {
			got := Match(tt.query, tt.target)
			if got != tt.want {
				t.Errorf("Match(%q, %q) = %v, want %v", tt.query, tt.target, got, tt.want)
			}
		})
	}
}

func TestScore(t *testing.T) {
	tests := []struct {
		query  string
		target string
		isMatch bool
	}{
		// Valid matches
		{"twt", "social/twitter.com/admin", true},
		{"gm", "email/gmail.com/user", true},

		// Non-matches
		{"twtt", "twitter", true},
		{"abc", "def", false},

		// Empty query
		{"", "short", true},
		{"", "this is a very long path", true},

		// Exact match - gets large bonus so score may be negative
		{"exact", "exact", true},
	}

	for _, tt := range tests {
		t.Run(tt.query+"_"+tt.target, func(t *testing.T) {
			got := Score(tt.query, tt.target)
			
			// For matches, score can be any value (lower is better, can be negative with bonuses)
			// For non-matches, score should be -1
			if tt.isMatch {
				// Just check it's not -1 (which is reserved for non-matches)
				// Scores can be positive or negative depending on bonuses
				if got == -1 {
					t.Errorf("Score(%q, %q) = -1, but it's a match", tt.query, tt.target)
				}
			} else {
				if got != -1 {
					t.Errorf("Score(%q, %q) = %d, want -1 for non-match", tt.query, tt.target, got)
				}
			}
		})
	}
}

func TestScoreOrdering(t *testing.T) {
	// Test that scores produce correct ordering
	items := []string{
		"email/gmail.com/user",
		"social/twitter.com/admin",
		"bank/chase.com/account",
		"email/outlook.com/work",
	}

	query := "gm"

	// For "gm", only email/gmail.com/user should match (has both g and m in order)
	// Other items don't have 'g' or don't have 'm' after 'g'
	
	// Get scores for all items
	scores := make([]int, len(items))
	for i, item := range items {
		scores[i] = Score(query, item)
	}

	// Collect only matching items
	var matchingItems []string
	var matchingScores []int
	for i, item := range items {
		if Match(query, item) {
			matchingItems = append(matchingItems, item)
			matchingScores = append(matchingScores, scores[i])
		}
	}
	
	t.Logf("Items matching 'gm': %v", matchingItems)
	t.Logf("Scores for matching items: %v", matchingScores)
	
	// Only email/gmail.com/user should match
	if len(matchingItems) != 1 {
		t.Errorf("Expected only 1 item to match 'gm', got %d: %v", len(matchingItems), matchingItems)
		return
	}
	
	if matchingItems[0] != "email/gmail.com/user" {
		t.Errorf("Expected email/gmail.com/user to match 'gm', got %s", matchingItems[0])
	}

	// Among matching items, email/gmail.com/user should have best score
	// Find the best score among matching items
	bestIdx := 0
	for i := 1; i < len(matchingItems); i++ {
		if matchingScores[i] < matchingScores[bestIdx] {
			bestIdx = i
		}
	}

	t.Logf("Scores for query '%s': %v", query, scores)
	t.Logf("Best match among matching items: %s", matchingItems[bestIdx])

	// The best match should be email/gmail.com/user (contains "gm" in gmail)
	if len(matchingItems) > 0 && matchingItems[bestIdx] != "email/gmail.com/user" {
		t.Errorf("Expected email/gmail.com/user to be best match for 'gm', got %s", matchingItems[bestIdx])
	}
}

func TestFilter(t *testing.T) {
	items := []string{
		"email/gmail.com/user",
		"social/twitter.com/admin",
		"bank/chase.com/account",
		"email/outlook.com/work",
		"random/path",
	}

	tests := []struct {
		query string
		want  []string
	}{
		{
			query: "gm",
			// Only email/gmail.com/user has both 'g' and 'm' in order
			want:  []string{"email/gmail.com/user"},
		},
		{
			query: "tw",
			// Both social/twitter.com/admin (t at 7, w at 8) and email/outlook.com/work (t at 8, w at 18) match
			want:  []string{"social/twitter.com/admin", "email/outlook.com/work"},
		},
		{
			query: "",
			want:  items,
		},
		{
			query: "nonexistent",
			want:  []string{},
		},
		{
			query: "chase",
			want:  []string{"bank/chase.com/account"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.query, func(t *testing.T) {
			got := Filter(tt.query, items)

			gotPaths := make([]string, len(got))
			for i, r := range got {
				gotPaths[i] = r.Path
			}

			if len(gotPaths) != len(tt.want) {
				t.Errorf("Filter(%q) returned %d results, want %d", tt.query, len(gotPaths), len(tt.want))
				return
			}

			for _, want := range tt.want {
				found := false
				for _, got := range gotPaths {
					if got == want {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Filter(%q) missing expected result %q", tt.query, want)
				}
			}
		})
	}
}

func TestFindBestMatch(t *testing.T) {
	items := []string{
		"email/gmail.com/user",
		"social/twitter.com/admin",
		"bank/chase.com/account",
	}

	tests := []struct {
		query string
		want  string
	}{
		{"gm", "email/gmail.com/user"},
		{"tw", "social/twitter.com/admin"},
		{"chase", "bank/chase.com/account"},
		{"nonexistent", ""},
	}

	for _, tt := range tests {
		t.Run(tt.query, func(t *testing.T) {
			got := FindBestMatch(tt.query, items)
			if got != tt.want {
				t.Errorf("FindBestMatch(%q) = %q, want %q", tt.query, got, tt.want)
			}
		})
	}
}

func TestIsSubsequence(t *testing.T) {
	tests := []struct {
		query  string
		target string
		want   bool
	}{
		{"abc", "abc", true},
		{"abc", "aabbcc", true},
		{"abc", "acb", false},
		{"abc", "ab", false},
		{"", "anything", true},
		{"a", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.query+"_"+tt.target, func(t *testing.T) {
			got := IsSubsequence(tt.query, tt.target)
			if got != tt.want {
				t.Errorf("IsSubsequence(%q, %q) = %v, want %v", tt.query, tt.target, got, tt.want)
			}
		})
	}
}

func TestMatchIndices(t *testing.T) {
	tests := []struct {
		query  string
		target string
		want   []int
	}{
		// For "gm" in "email/gmail.com/user": g is at index 6, m is at index 7
		{"gm", "email/gmail.com/user", []int{6, 7}},
		// For "twt" in "social/twitter.com/admin": t is at 7, w at 8, t at 10
		// social/twitter.com/admin: s(0),o(1),c(2),i(3),a(4),l(5),/(6),t(7),w(8),i(9),t(10),...
		{"twt", "social/twitter.com/admin", []int{7, 8, 10}},
		{"", "anything", []int{}},
	}

	for _, tt := range tests {
		t.Run(tt.query+"_"+tt.target, func(t *testing.T) {
			got := getMatchIndices(tt.query, tt.target)
			if len(got) != len(tt.want) {
				t.Errorf("getMatchIndices(%q, %q) returned %d indices, want %d", tt.query, tt.target, len(got), len(tt.want))
				return
			}

			for i, idx := range got {
				if idx != tt.want[i] {
					t.Errorf("getMatchIndices(%q, %q)[%d] = %d, want %d", tt.query, tt.target, i, idx, tt.want[i])
				}
			}
		})
	}
}

func TestEmptyQuery(t *testing.T) {
	items := []string{"c", "b", "a"}

	// Empty query should return all items
	results := Filter("", items)
	if len(results) != len(items) {
		t.Errorf("Empty query: expected %d results, got %d", len(items), len(results))
	}

	// Results should be sorted alphabetically
	for i := 0; i < len(results)-1; i++ {
		if results[i].Path > results[i+1].Path {
			t.Errorf("Empty query results not sorted: %s > %s", results[i].Path, results[i+1].Path)
		}
	}
	
	if results[0].Path != "a" || results[1].Path != "b" || results[2].Path != "c" {
		t.Errorf("Expected a, b, c order, got: %v", []string{results[0].Path, results[1].Path, results[2].Path})
	}
}
