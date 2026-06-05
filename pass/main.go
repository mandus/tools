// pass is a Windows-compatible replacement for the Unix password-store tool.
// It manages GPG-encrypted password files in ~/.password-store/ with git integration.
package main

import (
	"fmt"
	"os"

	"github.com/mandu/tools/pass/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
