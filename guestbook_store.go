package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

const guestbookEntriesFile = "guestbook-entries.json"

var (
	guestbookDataDir     = ".portico-data"
	guestbookSaveMu      sync.Mutex
)

func setGuestbookDataDir(dir string) {
	if dir != "" {
		guestbookDataDir = dir
	}
}

func guestbookEntriesPath() string {
	return filepath.Join(guestbookDataDir, guestbookEntriesFile)
}

func loadGuestbookEntries() ([]string, error) {
	path := guestbookEntriesPath()
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var out []string
	if err := json.Unmarshal(data, &out); err != nil {
		return nil, fmt.Errorf("guestbook entries: %w", err)
	}
	return out, nil
}

func saveGuestbookEntries(entries []string) error {
	guestbookSaveMu.Lock()
	defer guestbookSaveMu.Unlock()
	if err := os.MkdirAll(guestbookDataDir, 0o755); err != nil {
		return err
	}
	data, err := json.Marshal(entries)
	if err != nil {
		return err
	}
	path := guestbookEntriesPath()
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}
