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

func (d diffResult) isEmpty() bool {
	return len(d.Added) == 0 && len(d.Changed) == 0 && len(d.Deleted) == 0
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

// computeDiff builds the current local manifest, loads the last saved one,
// and returns what changed. Shared by both `push` and `diff` so the
// comparison logic lives in exactly one place.
func computeDiff(cfg Config) (diffResult, map[string]string, error) {
	newManifest, err := buildManifest(cfg.LocalDir)
	if err != nil {
		return diffResult{}, nil, err
	}

	oldManifest, err := loadPreviousManifest(manifestPath)
	if err != nil {
		return diffResult{}, nil, err
	}

	return diffManifests(oldManifest, newManifest), newManifest, nil
}
