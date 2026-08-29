package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/jlaffaye/ftp"
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

type Config struct {
	LocalDir  string `yaml:"local_dir"`
	RemoteDir string `yaml:"remote_dir"`
}

func loadConfig(path string) (Config, error) {
	var cfg Config

	data, err := os.ReadFile(path)
	if err != nil {
		return cfg, err
	}

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}

func buildManifest(dir string) (map[string]string, error) {
	manifest := make(map[string]string)

	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		h := sha256.New()
		if _, err := io.Copy(h, f); err != nil {
			return err
		}

		relPath, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}
		relPath = filepath.ToSlash(relPath)

		manifest[relPath] = hex.EncodeToString(h.Sum(nil))
		return nil
	})

	return manifest, err
}

func main() {
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

	conn, err := ftp.Dial(host)
	if err != nil {
		log.Fatalf("dial: %v", err)
	}
	defer conn.Quit()

	if err := conn.Login(user, pass); err != nil {
		log.Fatalf("login: %v", err)
	}

	manifest, err := buildManifest(cfg.LocalDir)
	if err != nil {
		log.Fatalf("build manifest: %v", err)
	}

	for relPath, hash := range manifest {
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

		fmt.Printf("uploaded %s (%s)\n", remotePath, hash[:8])
	}
}
