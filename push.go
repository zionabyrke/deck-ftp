package main

import (
	"os"
	"path/filepath"
	"time"

	"github.com/briandowns/spinner"
	"github.com/jlaffaye/ftp"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "Deploy changed files to the remote host",
	Run: func(cmd *cobra.Command, args []string) {
		if err := godotenv.Load(); err != nil {
			warn("no .env file found, falling back to real env vars")
		}

		cfg, err := loadConfig("deck.yaml")
		if err != nil {
			fail("load config: %v", err)
		}

		host := os.Getenv("DECK_FTP_HOST")
		user := os.Getenv("DECK_FTP_USER")
		pass := os.Getenv("DECK_FTP_PASS")

		if host == "" || user == "" || pass == "" {
			fail("missing DECK_FTP_HOST, DECK_FTP_USER, or DECK_FTP_PASS")
		}

		diff, newManifest, err := computeDiff(cfg)
		if err != nil {
			fail("compute diff: %v", err)
		}

		if diff.isEmpty() {
			warn("nothing to deploy")
			return
		}

		s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
		s.Suffix = " connecting..."
		s.Start()

		conn, err := ftp.Dial(host)
		if err != nil {
			s.Stop()
			fail("dial: %v", err)
		}
		defer conn.Quit()

		if err := conn.Login(user, pass); err != nil {
			s.Stop()
			fail("login: %v", err)
		}
		s.Stop()

		for _, relPath := range append(diff.Added, diff.Changed...) {
			localPath := filepath.Join(cfg.LocalDir, relPath)
			remotePath := cfg.RemoteDir + "/" + relPath

			s.Suffix = " uploading " + relPath
			s.Start()

			file, err := os.Open(localPath)
			if err != nil {
				s.Stop()
				fail("open %s: %v", localPath, err)
			}

			if err := conn.Stor(remotePath, file); err != nil {
				file.Close()
				s.Stop()
				fail("upload %s: %v", remotePath, err)
			}
			file.Close()
			s.Stop()

			success("uploaded %s", remotePath)
		}

		for _, relPath := range diff.Deleted {
			remotePath := cfg.RemoteDir + "/" + relPath
			if err := conn.Delete(remotePath); err != nil {
				fail("delete %s: %v", remotePath, err)
			}
			success("deleted %s", remotePath)
		}

		if err := saveManifest(manifestPath, newManifest); err != nil {
			fail("save manifest: %v", err)
		}

		success("done: %d added, %d changed, %d deleted",
			len(diff.Added), len(diff.Changed), len(diff.Deleted))
	},
}

func init() {
	rootCmd.AddCommand(pushCmd)
}
