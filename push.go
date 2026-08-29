package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/jlaffaye/ftp"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "Deploy changed files to the remote host",
	Run: func(cmd *cobra.Command, args []string) {
		if err := godotenv.Load(); err != nil {
			log.Println("no .env file found, falling back to real env vars")
		}

		cfg, err := loadConfig("deck.yaml")
		if err != nil {
			log.Fatalf("load config: %v", err)
		}

		host := os.Getenv("DECK_FTP_HOST")
		user := os.Getenv("DECK_FTP_USER")
		pass := os.Getenv("DECK_FTP_PASS")

		if host == "" || user == "" || pass == "" {
			log.Fatal("missing DECK_FTP_HOST, DECK_FTP_USER, or DECK_FTP_PASS")
		}

		diff, newManifest, err := computeDiff(cfg)
		if err != nil {
			log.Fatalf("compute diff: %v", err)
		}

		if diff.isEmpty() {
			fmt.Println("nothing to deploy")
			return
		}

		conn, err := ftp.Dial(host)
		if err != nil {
			log.Fatalf("dial: %v", err)
		}
		defer conn.Quit()

		if err := conn.Login(user, pass); err != nil {
			log.Fatalf("login: %v", err)
		}

		for _, relPath := range append(diff.Added, diff.Changed...) {
			localPath := filepath.Join(cfg.LocalDir, relPath)
			remotePath := cfg.RemoteDir + "/" + relPath

			file, err := os.Open(localPath)
			if err != nil {
				log.Fatalf("open %s: %v", localPath, err)
			}

			if err := conn.Stor(remotePath, file); err != nil {
				file.Close()
				log.Fatalf("upload %s: %v", remotePath, err)
			}
			file.Close()

			fmt.Println("uploaded", remotePath)
		}

		for _, relPath := range diff.Deleted {
			remotePath := cfg.RemoteDir + "/" + relPath
			if err := conn.Delete(remotePath); err != nil {
				log.Fatalf("delete %s: %v", remotePath, err)
			}
			fmt.Println("deleted", remotePath)
		}

		if err := saveManifest(manifestPath, newManifest); err != nil {
			log.Fatalf("save manifest: %v", err)
		}

		fmt.Printf("done: %d added, %d changed, %d deleted\n",
			len(diff.Added), len(diff.Changed), len(diff.Deleted))
	},
}

func init() {
	rootCmd.AddCommand(pushCmd)
}
