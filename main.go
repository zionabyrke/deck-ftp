package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jlaffaye/ftp"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, falling back to real env vars")
	}

	host := os.Getenv("DECK_FTP_HOST")
	// ... rest stays exactly the same
	user := os.Getenv("DECK_FTP_USER")
	pass := os.Getenv("DECK_FTP_PASS")

	if host == "" || user == "" || pass == "" {
		log.Fatal("missing DECK_FTP_HOST, DECK_FTP_USER, or DECK_FTP_PASS")
	}

	localPath := "spike-test.txt"
	if err := os.WriteFile(localPath, []byte("hello from DECK\n"), 0644); err != nil {
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

	file, err := os.Open(localPath)
	if err != nil {
		log.Fatalf("open local file: %v", err)
	}
	defer file.Close()

	remotePath := "htdocs/spike-test.txt"
	if err := conn.Stor(remotePath, file); err != nil {
		log.Fatalf("upload: %v", err)
	}

	fmt.Println("uploaded", remotePath)
}
