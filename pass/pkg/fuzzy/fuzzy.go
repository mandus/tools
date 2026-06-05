// Package fuzzy provides fuzzy string matching functionality for the pass tool.
// It implements subsequence matching (characters must appear in order but not consecutively)
// and scoring for ranking matches.
package fuzzy

import (
	"strings"
)

// MatchResult represents a fuzzy match result with score and match positions.
type MatchResult struct {
	Path       string
	Score      int
	MatchIndices []int // Indices of matching characters in Path (byte indices)
}

// Match checks if query is a subsequence of target (case-insensitive).
// Returns true if all characters in query appear in target in the same order.
func Match(query, target string) bool {
	if query == "" {
		return true // Empty query matches everything
	}
	if target == "" {
		return false
	}

	lowerQuery := strings.ToLower(query)
	lowerTarget := strings.ToLower(target)

	queryIdx := 0
	for i := 0; i < len(lowerTarget) && queryIdx < len(lowerQuery); i++ {
		if lowerTarget[i] == lowerQuery[queryIdx] {
			queryIdx++
		}
	}

	return queryIdx == len(lowerQuery)
}

// Score calculates the match quality between query and target.
// Lower score means better match.
// Scoring factors:
// - Characters must match in order (subsequence)
// - Earlier matches are better
// - Consecutive matches get bonus
// - Matches at start of path components get bonus
// - Shorter paths are better
func Score(query, target string) int {
	if query == "" {
		// Empty query: score based on path length (shorter is better)
		return len(target)
	}

	lowerQuery := strings.ToLower(query)
	lowerTarget := strings.ToLower(target)

	// First, check if it's a valid match
	if !Match(query, target) {
		return -1 // No match
	}

	score := 0
	queryIdx := 0
	prevMatch := false
	pathComponentStart := true

	for i := 0; i < len(lowerTarget); i++ {
		char := lowerTarget[i]
		isMatch := queryIdx < len(lowerQuery) && char == lowerQuery[queryIdx]

		if isMatch {
			// Bonus for match at start of path component
			if pathComponentStart {
				score += 50
			}

			// Bonus for consecutive matches
			if prevMatch {
				score += 20
			}

			// Main score: earlier matches are much better
			// Use a decaying bonus based on position
			score += 1000 - (i * 2)

			queryIdx++
			prevMatch = true
		} else {
			// Small penalty for non-matching characters
			score += 1
			prevMatch = false
		}

		// Track path component starts (after /)
		if char == '/' {
			pathComponentStart = true
		} else {
			pathComponentStart = false
		}
	}

	// Penalty for long paths (encourages shorter matches)
	score -= len(target) * 5

	// Bonus for exact match
	if strings.ToLower(query) == strings.ToLower(target) {
		score -= 10000 // Large bonus for exact match
	}

	return score
}

// Filter filters a list of items using fuzzy matching and returns sorted results.
// Results are sorted by score (best matches first), then alphabetically for ties.
func Filter(query string, items []string) []MatchResult {
	if query == "" {
		// Empty query: return all items sorted alphabetically
		results := make([]MatchResult, len(items))
		for i, item := range items {
			results[i] = MatchResult{
				Path:       item,
				Score:      len(item),
				MatchIndices: []int{},
			}
		}
		// Sort alphabetically for empty query
		for i := 0; i < len(results)-1; i++ {
			for j := 0; j < len(results)-i-1; j++ {
				if results[j].Path > results[j+1].Path {
					results[j], results[j+1] = results[j+1], results[j]
				}
			}
		}
		return results
	}

	var results []MatchResult
	for _, item := range items {
		if Match(query, item) {
			matchIndices := getMatchIndices(query, item)
			results = append(results, MatchResult{
				Path:       item,
				Score:      Score(query, item),
				MatchIndices: matchIndices,
			})
		}
	}

	// Sort by score (lower is better), then by path for ties
	for i := 0; i < len(results)-1; i++ {
		for j := 0; j < len(results)-i-1; j++ {
			if results[j].Score > results[j+1].Score {
				results[j], results[j+1] = results[j+1], results[j]
			} else if results[j].Score == results[j+1].Score {
				// Tie-breaker: alphabetical order
				if results[j].Path > results[j+1].Path {
					results[j], results[j+1] = results[j+1], results[j]
				}
			}
		}
	}

	return results
}

// getMatchIndices returns the byte indices of matching characters in target.
func getMatchIndices(query, target string) []int {
	if query == "" {
		return []int{}
	}

	lowerQuery := strings.ToLower(query)
	lowerTarget := strings.ToLower(target)

	var indices []int
	queryIdx := 0

	for i := 0; i < len(lowerTarget); i++ {
		if queryIdx < len(lowerQuery) && lowerTarget[i] == lowerQuery[queryIdx] {
			indices = append(indices, i)
			queryIdx++
		}
	}

	return indices
}

// FindBestMatch returns the best matching item from a list, or empty string if no match.
func FindBestMatch(query string, items []string) string {
	results := Filter(query, items)
	if len(results) > 0 {
		return results[0].Path
	}
	return ""
}

// IsSubsequence checks if query is a subsequence of target (case-sensitive version).
func IsSubsequence(query, target string) bool {
	if query == "" {
		return true
	}
	queryIdx := 0
	for i := 0; i < len(target) && queryIdx < len(query); i++ {
		if target[i] == query[queryIdx] {
			queryIdx++
		}
	}
	return queryIdx == len(query)
}
