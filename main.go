package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jlaffaye/ftp"
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

type Config struct {
	LocalFile  string `yaml:"local_file"`
	RemoteFile string `yaml:"remote_file"`
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

	if err := os.WriteFile(cfg.LocalFile, []byte("hello from DECK\n"), 0644); err != nil {
		log.Fatalf("write local file: %v", err)
	}

	conn, err := ftp.Dial(host)
	if err != nil {
		log.Fatalf("dial: %v", err)
	}
	defer conn.Quit()

	if err := conn.Login(user, pass); err != nil {
		log.Fatalf("login: %v", err)
	}

	file, err := os.Open(cfg.LocalFile)
	if err != nil {
		log.Fatalf("open local file: %v", err)
	}
	defer file.Close()

	if err := conn.Stor(cfg.RemoteFile, file); err != nil {
		log.Fatalf("upload: %v", err)
	}

	fmt.Println("uploaded", cfg.RemoteFile)
}
