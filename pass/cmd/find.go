package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mandu/tools/pass/cmd/tree"
	"github.com/mandu/tools/pass/pkg/filesystem"
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
	ignoreCaseFlag bool
	flatFlag       bool
	noTreeFlag     bool
)

func addFindCmd() {
	findCmd.Flags().BoolVarP(&ignoreCaseFlag, "ignore-case", "i", false, "Case-insensitive search")
	findCmd.Flags().BoolVarP(&flatFlag, "flat", "f", false, "Output flat list instead of tree")
	findCmd.Flags().BoolVar(&noTreeFlag, "no-tree", false, "Output flat list instead of tree")
	rootCmd.AddCommand(findCmd)
}

// renderTreeNode renders a tree node and its children with proper box-drawing characters.
// prefix is the indentation/vertical connector prefix for this node.
// isLast indicates whether this node is the last child of its parent.
func renderTreeNode(node *tree.TreeNode, prefix string, isLast bool) {
	// Determine the connector for this node
	connector := "\u2514\u2500\u2500 " // └──
	if !isLast {
		connector = "\u251C\u2500\u2500 " // ├──
	}

	// Format the node name
	name := node.Name
	if node.IsDir {
		name += "/"
	}

	// Print this node
	fmt.Print(prefix + connector + name + "\n")

	// Build child prefix
	childPrefix := prefix
	if isLast {
		childPrefix += "    "
	} else {
		childPrefix += "\u2502   " // │   
	}

	// Render children
	for i, child := range node.Children {
		childIsLast := i == len(node.Children)-1
		renderTreeNode(child, childPrefix, childIsLast)
	}
}

// findPasswords searches for passwords containing the search string
func findPasswords(searchString string, flat bool) error {
	if searchString == "" {
		return fmt.Errorf("pass: search string cannot be empty")
	}

	storeDir := GetPasswordStoreDir()

	// Check if store exists
	if _, err := os.Stat(storeDir); os.IsNotExist(err) {
		fmt.Println() // Print nothing if store doesn't exist
		return nil
	}

	// Prepare search string
	target := searchString
	if ignoreCaseFlag {
		target = strings.ToLower(searchString)
	}

	// Walk the directory tree
	var results []string
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
			// Get relative path from store directory
			relPath, err := filepath.Rel(storeDir, path)
			if err != nil {
				return err
			}

			// Normalize path separators
			relPath = filesystem.NormalizePathForDisplay(relPath)

			// Strip .gpg extension
			passwordPath := strings.TrimSuffix(relPath, ".gpg")

			// Check if path contains search string
			checkPath := passwordPath
			if ignoreCaseFlag {
				checkPath = strings.ToLower(passwordPath)
			}

			if strings.Contains(checkPath, target) {
				results = append(results, passwordPath)
			}
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("pass: failed to walk directory: %v", err)
	}

	// Sort results
	for i := 0; i < len(results)-1; i++ {
		for j := 0; j < len(results)-i-1; j++ {
			if results[j] > results[j+1] {
				results[j], results[j+1] = results[j+1], results[j]
			}
		}
	}

	// Print results
	if flat {
		// Original flat list output
		for _, result := range results {
			fmt.Println(result)
		}
	} else {
		// Tree view output
		if len(results) > 0 {
			treeRoot := tree.BuildTreeFromPaths(results)
			// Render each top-level child with empty prefix
			for i := range treeRoot.Children {
				isLast := i == len(treeRoot.Children)-1
				renderTreeNode(treeRoot.Children[i], "", isLast)
			}
		}
	}

	return nil
}
