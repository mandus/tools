package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mandu/tools/pass/pkg/filesystem"
	"github.com/mandu/tools/pass/pkg/fuzzy"
	"github.com/mandu/tools/pass/pkg/gpg"
	"github.com/mandu/tools/pass/pkg/terminal"
)

// FuzzySearchMode represents the mode of fuzzy search (show, clip, rm)
type FuzzySearchMode int

const (
	// FuzzyModeShow - show the selected password
	FuzzyModeShow FuzzySearchMode = iota
	// FuzzyModeClip - copy the selected password to clipboard
	FuzzyModeClip
	// FuzzyModeRm - delete the selected password
	FuzzyModeRm
)

// FuzzySearchState holds the state of the fuzzy search UI
type FuzzySearchState struct {
	passwords    []string
	filtered    []fuzzy.MatchResult
	query       string
	queryCursor int
	selectedIdx int
	scrollPos   int
	width       int
	height      int
	listHeight  int
	mode        FuzzySearchMode
	keyReader   *terminal.KeyReader
}

func newFuzzySearchState(passwords []string, mode FuzzySearchMode) *FuzzySearchState {
	w, h, _ := terminal.GetSize()
	listHeight := h - 2
	if listHeight < 1 {
		listHeight = 1
	}
	return &FuzzySearchState{
		passwords:    passwords,
		filtered:    nil,
		query:       "",
		queryCursor: 0,
		selectedIdx: 0,
		scrollPos:   0,
		width:       w,
		height:      h,
		listHeight:  listHeight,
		mode:        mode,
	}
}

func (s *FuzzySearchState) updateFilter() {
	s.filtered = fuzzy.Filter(s.query, s.passwords)
	s.selectedIdx = 0
	s.scrollPos = 0
	if len(s.filtered) > 0 && s.selectedIdx >= len(s.filtered) {
		s.selectedIdx = len(s.filtered) - 1
	}
}

func (s *FuzzySearchState) getVisibleItems() []fuzzy.MatchResult {
	end := s.scrollPos + s.listHeight
	if end > len(s.filtered) {
		end = len(s.filtered)
	}
	return s.filtered[s.scrollPos:end]
}

func (s *FuzzySearchState) scrollToItem(idx int) {
	if idx < s.scrollPos {
		s.scrollPos = idx
	} else if idx >= s.scrollPos+s.listHeight {
		s.scrollPos = idx - s.listHeight + 1
	}
	if s.scrollPos < 0 {
		s.scrollPos = 0
	}
	maxScroll := 0
	if len(s.filtered) > s.listHeight {
		maxScroll = len(s.filtered) - s.listHeight
	}
	if s.scrollPos > maxScroll {
		s.scrollPos = maxScroll
	}
}

func (s *FuzzySearchState) handleKey(key terminal.Key) (bool, string, error) {
	// Exit keys
	if key.IsEscape || (key.IsCtrl && key.CtrlChar == 'C') || (key.IsCtrl && key.CtrlChar == 'D') {
		return false, "", nil
	}

	// Enter - select
	if key.IsEnter {
		if len(s.filtered) == 0 {
			return true, "", nil
		}
		return false, s.filtered[s.selectedIdx].Path, nil
	}

	// Arrow keys
	if key.IsArrow {
		switch key.ArrowDir {
		case "up":
			if s.selectedIdx > 0 {
				s.selectedIdx--
				s.scrollToItem(s.selectedIdx)
			}
		case "down":
			if s.selectedIdx < len(s.filtered)-1 {
				s.selectedIdx++
				s.scrollToItem(s.selectedIdx)
			}
		case "left":
			if s.queryCursor > 0 {
				s.queryCursor--
			}
		case "right":
			if s.queryCursor < len(s.query) {
				s.queryCursor++
			}
		}
		return true, "", nil
	}

	// Home - start of query
	if key.IsHome || (key.IsCtrl && key.CtrlChar == 'A') {
		s.queryCursor = 0
		return true, "", nil
	}

	// End - end of query
	if key.IsEnd || (key.IsCtrl && key.CtrlChar == 'E') {
		s.queryCursor = len(s.query)
		return true, "", nil
	}

	// Ctrl+K - clear from cursor to end
	if key.IsCtrl && key.CtrlChar == 'K' {
		if s.queryCursor < len(s.query) {
			s.query = s.query[:s.queryCursor]
		}
		s.updateFilter()
		return true, "", nil
	}

	// Backspace
	if key.IsBackspace || key.Rune == 8 || key.Rune == 127 {
		if s.queryCursor > 0 {
			s.query = s.query[:s.queryCursor-1] + s.query[s.queryCursor:]
			s.queryCursor--
			s.updateFilter()
		}
		return true, "", nil
	}

	// Delete
	if key.IsDelete {
		if s.queryCursor < len(s.query) {
			s.query = s.query[:s.queryCursor] + s.query[s.queryCursor+1:]
			s.updateFilter()
		}
		return true, "", nil
	}

	// Ctrl+L - clear query
	if key.IsCtrl && key.CtrlChar == 'L' {
		s.query = ""
		s.queryCursor = 0
		s.updateFilter()
		return true, "", nil
	}

	// Ctrl+W - delete word
	if key.IsCtrl && key.CtrlChar == 'W' {
		start := s.queryCursor
		for i := s.queryCursor - 1; i >= 0; i-- {
			if s.query[i] == ' ' || s.query[i] == '/' {
				start = i + 1
				break
			}
		}
		if start < s.queryCursor {
			s.query = s.query[:start] + s.query[s.queryCursor:]
			s.queryCursor = start
			s.updateFilter()
		}
		return true, "", nil
	}

	// Tab - cycle through results
	if key.IsTab {
		if len(s.filtered) > 1 {
			s.selectedIdx = (s.selectedIdx + 1) % len(s.filtered)
			s.scrollToItem(s.selectedIdx)
		}
		return true, "", nil
	}

	// Page Up
	if key.IsPageUp {
		s.selectedIdx -= s.listHeight
		if s.selectedIdx < 0 {
			s.selectedIdx = 0
		}
		s.scrollToItem(s.selectedIdx)
		return true, "", nil
	}

	// Page Down
	if key.IsPageDown {
		s.selectedIdx += s.listHeight
		if s.selectedIdx >= len(s.filtered) {
			s.selectedIdx = len(s.filtered) - 1
		}
		if s.selectedIdx >= 0 {
			s.scrollToItem(s.selectedIdx)
		}
		return true, "", nil
	}

	// Printable characters
	if key.Rune > 31 && key.Rune != 127 {
		s.query = s.query[:s.queryCursor] + string(key.Rune) + s.query[s.queryCursor:]
		s.queryCursor++
		s.updateFilter()
	}

	return true, "", nil
}

func (s *FuzzySearchState) getPrompt() string {
	switch s.mode {
	case FuzzyModeClip:
		return "Copy: "
	case FuzzyModeRm:
		return "Remove: "
	default:
		return "Search: "
	}
}

func (s *FuzzySearchState) render() {
	terminal.ClearScreenFunc()

	// Header
	var header string
	switch s.mode {
	case FuzzyModeClip:
		header = "Select password to copy (Enter to copy, Esc to cancel):"
	case FuzzyModeRm:
		header = "Select password to remove (Enter to delete, Esc to cancel):"
	default:
		header = "Select password (Enter to show, Esc to cancel):"
	}
	fmt.Println(terminal.Truncate(header, s.width))

	// List
	visible := s.getVisibleItems()
	for i, item := range visible {
		prefix := "  "
		actualIdx := s.scrollPos + i
		if actualIdx == s.selectedIdx {
			prefix = "> "
		}
		path := terminal.HighlightMatch(item.Path, item.MatchIndices)
		line := prefix + terminal.Truncate(path, s.width-2)
		fmt.Println(terminal.PadRight(line, s.width))
	}

	// Fill empty space
	for len(visible) < s.listHeight {
		fmt.Println(terminal.PadRight("", s.width))
	}

	// Prompt
	prompt := s.getPrompt() + s.query
	fmt.Print(terminal.PadRight(prompt, s.width))

	// Position cursor
	if terminal.SupportsANSI() {
		col := len(s.getPrompt()) + s.queryCursor
		fmt.Printf("\033[%dA\033[%dC", len(visible)+1, col)
	}
}

// InteractiveFuzzySearch runs the interactive fuzzy search
func InteractiveFuzzySearch(mode FuzzySearchMode) (string, error) {
	storeDir := GetPasswordStoreDir()
	if _, err := os.Stat(storeDir); os.IsNotExist(err) {
		return "", fmt.Errorf("pass: password store does not exist. Use 'pass insert' to create it.")
	}

	passwords, err := collectAllPasswords(storeDir)
	if err != nil {
		return "", err
	}
	if len(passwords) == 0 {
		return "", fmt.Errorf("pass: no passwords found in store")
	}

	state := newFuzzySearchState(passwords, mode)
	state.updateFilter()

	keyReader, err := terminal.NewKeyReader()
	if err != nil {
		return "", fmt.Errorf("pass: %v", err)
	}
	defer keyReader.Close()

	if err := keyReader.EnableRawMode(); err != nil {
		fmt.Fprintf(os.Stderr, "warning: raw mode not available: %v\n", err)
	}

	if terminal.SupportsANSI() {
		defer terminal.ShowCursorFunc()
		terminal.HideCursorFunc()
	}

	state.render()

	for {
		key, err := keyReader.ReadKey()
		if err != nil {
			break
		}

		cont, selected, err := state.handleKey(key)
		if err != nil {
			fmt.Println()
			return "", err
		}
		if !cont {
			fmt.Println()
			return selected, nil
		}
		state.render()
	}

	fmt.Println()
	return "", nil
}

// RunInteractiveFuzzySearch runs interactive search and performs the action
func RunInteractiveFuzzySearch(mode FuzzySearchMode) error {
	selected, err := InteractiveFuzzySearch(mode)
	if err != nil {
		return err
	}
	if selected == "" {
		return nil
	}

	switch mode {
	case FuzzyModeShow:
		return showPassword(selected)
	case FuzzyModeClip:
		fullPath := filesystem.JoinPath(GetPasswordStoreDir(), filesystem.NormalizePath(selected)+".gpg")
		password, err := gpg.DecryptFile(fullPath)
		if err != nil {
			return err
		}
		if err := filesystem.CopyToClipboard(password); err != nil {
			return fmt.Errorf("pass: failed to copy to clipboard: %v", err)
		}
		fmt.Printf("Copied %s to clipboard.\n", selected)
		go filesystem.StartClipboardClearTimer()
		return nil
	case FuzzyModeRm:
		return removePassword(selected, false, false)
	}
	return nil
}

// fuzzySearchMode is the entry point called from root.go
func fuzzySearchMode() error {
	return RunInteractiveFuzzySearch(FuzzyModeShow)
}

// collectAllPasswords walks the password store and collects all password paths.
func collectAllPasswords(storeDir string) ([]string, error) {
	var passwords []string

	err := filepath.Walk(storeDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip .git directory
		if info.IsDir() && info.Name() == ".git" {
			return filepath.SkipDir
		}

		// Only process .gpg files
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".gpg") {
			relPath, err := filepath.Rel(storeDir, path)
			if err != nil {
				return err
			}
			relPath = filesystem.NormalizePathForDisplay(relPath)
			passwordPath := strings.TrimSuffix(relPath, ".gpg")
			passwords = append(passwords, passwordPath)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("pass: failed to list passwords: %v", err)
	}

	return passwords, nil
}

// showSelectedPassword shows the password at the selected path.
func showSelectedPassword(path string) error {
	return showPassword(path)
}

// UseFuzzyPath finds the best fuzzy match for a query
func UseFuzzyPath(query string) (string, error) {
	storeDir := GetPasswordStoreDir()

	passwords, err := collectAllPasswords(storeDir)
	if err != nil {
		return "", err
	}

	if len(passwords) == 0 {
		return "", fmt.Errorf("pass: no passwords found in store")
	}

	bestMatch := fuzzy.FindBestMatch(query, passwords)
	if bestMatch == "" {
		return "", fmt.Errorf("pass: no match found for %q", query)
	}

	return bestMatch, nil
}
