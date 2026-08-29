package main

import (
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show current config and whether a deploy is pending",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := loadConfig("deck.yaml")
		if err != nil {
			fail("load config: %v", err)
		}

		successColor.Printf("local_dir:  %s\n", cfg.LocalDir)
		successColor.Printf("remote_dir: %s\n", cfg.RemoteDir)

		diff, _, err := computeDiff(cfg)
		if err != nil {
			fail("compute diff: %v", err)
		}

		if diff.isEmpty() {
			success("status:     up to date")
			return
		}

		warn("status:     %d added, %d changed, %d deleted (run 'deck-ftp push' to deploy)",
			len(diff.Added), len(diff.Changed), len(diff.Deleted))
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
