package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

type WatchedDB struct {
	Episodes map[string]WatchedEntry `json:"episodes"`
}

type WatchedEntry struct {
	WatchedAt time.Time `json:"watched_at"`
	Title     string    `json:"title"`
}

func watchedPath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "xtreamgo", "watched.json"), nil
}

func loadWatched() (*WatchedDB, error) {
	path, err := watchedPath()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return &WatchedDB{Episodes: make(map[string]WatchedEntry)}, nil
	}
	if err != nil {
		return nil, err
	}
	var db WatchedDB
	if err := json.Unmarshal(data, &db); err != nil {
		return nil, err
	}
	if db.Episodes == nil {
		db.Episodes = make(map[string]WatchedEntry)
	}
	return &db, nil
}

func saveWatched(db *WatchedDB) error {
	path, err := watchedPath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(db, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}

func (db *WatchedDB) IsWatched(episodeID string) bool {
	_, ok := db.Episodes[episodeID]
	return ok
}

func (db *WatchedDB) Mark(episodeID, title string) {
	db.Episodes[episodeID] = WatchedEntry{
		WatchedAt: time.Now(),
		Title:     title,
	}
}

func (db *WatchedDB) Unmark(episodeID string) {
	delete(db.Episodes, episodeID)
}

func (db *WatchedDB) Toggle(episodeID, title string) bool {
	if db.IsWatched(episodeID) {
		db.Unmark(episodeID)
		return false
	}
	db.Mark(episodeID, title)
	return true
}