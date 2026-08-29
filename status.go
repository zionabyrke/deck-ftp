package main

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show current config and whether a deploy is pending",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := loadConfig("deck.yaml")
		if err != nil {
			log.Fatalf("load config: %v", err)
		}

		fmt.Println("local_dir:  " + cfg.LocalDir)
		fmt.Println("remote_dir: " + cfg.RemoteDir)

		diff, _, err := computeDiff(cfg)
		if err != nil {
			log.Fatalf("compute diff: %v", err)
		}

		if diff.isEmpty() {
			fmt.Println("status:     up to date")
			return
		}

		fmt.Printf("status:     %d added, %d changed, %d deleted (run 'deck-ftp push' to deploy)\n",
			len(diff.Added), len(diff.Changed), len(diff.Deleted))
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
