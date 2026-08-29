package main

import (
	"os"

	"github.com/spf13/cobra"
)

const defaultConfig = `local_dir: site
remote_dir: htdocs
`

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Create a starter deck.yaml in the current directory",
	Run: func(cmd *cobra.Command, args []string) {
		if _, err := os.Stat("deck.yaml"); err == nil {
			fail("deck.yaml already exists, refusing to overwrite")
		}

		if err := os.WriteFile("deck.yaml", []byte(defaultConfig), 0644); err != nil {
			fail("write deck.yaml: %v", err)
		}

		success("created deck.yaml")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
