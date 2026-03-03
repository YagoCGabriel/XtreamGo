package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	Servers []Server `json:"servers"`
	Current int      `json:"current"`
}

type Server struct {
	Name       string `json:"name"`
	URL        string `json:"url"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	Player     string `json:"player"`
	HWDec      string `json:"hwdec"`
	Fullscreen bool   `json:"fullscreen"`
}

func configPath() (string, error) {
	exePath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("nao foi possivel determinar o diretorio do executavel: %w", err)
	}
	// filepath.EvalSymlinks resolve casos onde o binario e um symlink
	exePath, err = filepath.EvalSymlinks(exePath)
	if err != nil {
		return "", err
	}
	return filepath.Join(filepath.Dir(exePath), "config.json"), nil
}

func loadConfig() (*Config, error) {
	path, err := configPath()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return &Config{}, nil
	}
	if err != nil {
		return nil, err
	}
	var cfg Config
	return &cfg, json.Unmarshal(data, &cfg)
}

func saveConfig(cfg *Config) error {
	path, err := configPath()
	if err != nil {
		return err
	}
	// garante que o diretorio existe (util se rodar de symlink em outro dir)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}

func (c *Config) activeServer() (*Server, error) {
	if len(c.Servers) == 0 {
		return nil, fmt.Errorf("nenhum servidor — execute: xtream-mpv add")
	}
	if c.Current >= len(c.Servers) {
		c.Current = 0
	}
	return &c.Servers[c.Current], nil
}