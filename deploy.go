package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	LocalDir  string `yaml:"local_dir"`
	RemoteDir string `yaml:"remote_dir"`
}

const manifestPath = ".deck/manifest.json"

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

func loadPreviousManifest(path string) (map[string]string, error) {
	manifest := make(map[string]string)

	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return manifest, nil
	}
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, &manifest); err != nil {
		return nil, err
	}

	return manifest, nil
}

func saveManifest(path string, manifest map[string]string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

type diffResult struct {
	Added   []string
	Changed []string
	Deleted []string
}

func diffManifests(oldM, newM map[string]string) diffResult {
	var d diffResult

	for path, hash := range newM {
		oldHash, existed := oldM[path]
		if !existed {
			d.Added = append(d.Added, path)
		} else if oldHash != hash {
			d.Changed = append(d.Changed, path)
		}
	}

	for path := range oldM {
		if _, stillExists := newM[path]; !stillExists {
			d.Deleted = append(d.Deleted, path)
		}
	}

	return d
}
