package main

import (
	"github.com/spf13/cobra"
)

var diffCmd = &cobra.Command{
	Use:   "diff",
	Short: "Show what would change on the next push, without deploying",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := loadConfig("deck.yaml")
		if err != nil {
			fail("load config: %v", err)
		}

		diff, _, err := computeDiff(cfg)
		if err != nil {
			fail("compute diff: %v", err)
		}

		if diff.isEmpty() {
			warn("no changes")
			return
		}

		for _, path := range diff.Added {
			success("+ %s", path)
		}
		for _, path := range diff.Changed {
			warn("~ %s", path)
		}
		for _, path := range diff.Deleted {
			errorColor.Printf("- %s\n", path)
		}
	},
}

func init() {
	rootCmd.AddCommand(diffCmd)
}
