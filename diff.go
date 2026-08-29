package main

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var diffCmd = &cobra.Command{
	Use:   "diff",
	Short: "Show what would change on the next push, without deploying",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := loadConfig("deck.yaml")
		if err != nil {
			log.Fatalf("load config: %v", err)
		}

		diff, _, err := computeDiff(cfg)
		if err != nil {
			log.Fatalf("compute diff: %v", err)
		}

		if diff.isEmpty() {
			fmt.Println("no changes")
			return
		}

		for _, path := range diff.Added {
			fmt.Println("+ " + path)
		}
		for _, path := range diff.Changed {
			fmt.Println("~ " + path)
		}
		for _, path := range diff.Deleted {
			fmt.Println("- " + path)
		}
	},
}

func init() {
	rootCmd.AddCommand(diffCmd)
}
