package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

const hookScript = `#!/bin/sh
branch=$(git rev-parse --abbrev-ref HEAD)
if [ "$branch" = "main" ]; then
  ./deck-ftp push
fi
`

var installHookCmd = &cobra.Command{
	Use:   "install-hook",
	Short: "Install a git post-commit hook that auto-deploys on main",
	Run: func(cmd *cobra.Command, args []string) {
		hookPath := filepath.Join(".git", "hooks", "post-commit")

		if _, err := os.Stat(".git"); os.IsNotExist(err) {
			log.Fatal("no .git directory found, run this from inside a git repo")
		}

		if err := os.WriteFile(hookPath, []byte(hookScript), 0755); err != nil {
			log.Fatalf("write hook: %v", err)
		}

		fmt.Println("installed post-commit hook at", hookPath)
	},
}

func init() {
	rootCmd.AddCommand(installHookCmd)
}
