package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	Servers  []Server `json:"servers"`
	Current  int      `json:"current"`
	Language string   `json:"language"` // "pt-BR" ou "en-US"
}

type Server struct {
	Name       string `json:"name"`
	URL        string `json:"url"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	Player     string `json:"player"`     // "mpv", "vlc" ou "kmplayer"
	HWDec      string `json:"hwdec"`      // "no", "auto", "vaapi", "nvdec", etc.
	Fullscreen bool   `json:"fullscreen"`
}

func configPath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "xtream-mpv", "config.json"), nil
}

func loadConfig() (*Config, error) {
	path, err := configPath()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return &Config{Language: "pt-BR"}, nil
	}
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	if cfg.Language == "" {
		cfg.Language = "pt-BR"
	}
	return &cfg, nil
}

func saveConfig(cfg *Config) error {
	path, err := configPath()
	if err != nil {
		return err
	}
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
		// mensagem bilíngue pois o idioma pode não estar carregado ainda
		return nil, fmt.Errorf("no server configured — run: xtreamgo add\nnenhum servidor — execute: xtreamgo add")
	}
	if c.Current >= len(c.Servers) {
		c.Current = 0
	}
	return &c.Servers[c.Current], nil
}
