package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Profile struct {
	Name    string `json:"name"`
	Token   string `json:"token"`
	ChatID  string `json:"chat_id"`
	Default bool   `json:"default,omitempty"`
}

type Config struct {
	Profiles []Profile `json:"profiles"`
}

func ConfigPath() (string, error) {
	dir, err := configDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "config.json"), nil
}

func configDir() (string, error) {
	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		return filepath.Join(xdg, "auto-message"), nil
	}
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("cannot determine config directory: %w", err)
	}
	return filepath.Join(dir, "auto-message"), nil
}

func Load() (*Config, error) {
	path, err := ConfigPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &Config{}, nil
		}
		return nil, fmt.Errorf("cannot read config file: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("cannot parse config file: %w", err)
	}
	return &cfg, nil
}

func Save(cfg *Config) error {
	path, err := ConfigPath()
	if err != nil {
		return err
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("cannot create config directory: %w", err)
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("cannot serialize config: %w", err)
	}

	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("cannot write config file: %w", err)
	}
	return nil
}

func GetDefault(cfg *Config) (*Profile, error) {
	for i := range cfg.Profiles {
		if cfg.Profiles[i].Default {
			return &cfg.Profiles[i], nil
		}
	}
	if len(cfg.Profiles) > 0 {
		return &cfg.Profiles[0], nil
	}
	return nil, fmt.Errorf("no config profiles found. Run 'auto-message config add' to create one, or set AUTO_MESSAGE_TOKEN and AUTO_MESSAGE_CHAT_ID environment variables")
}

func GetProfile(cfg *Config, name string) (*Profile, error) {
	for i := range cfg.Profiles {
		if cfg.Profiles[i].Name == name {
			return &cfg.Profiles[i], nil
		}
	}
	return nil, fmt.Errorf("config '%s' not found. Run 'auto-message config list' to see available profiles", name)
}

func AddProfile(cfg *Config, p Profile) error {
	for _, existing := range cfg.Profiles {
		if existing.Name == p.Name {
			return fmt.Errorf("a profile named '%s' already exists. Use a different name or remove the existing one first", p.Name)
		}
	}

	if p.Default {
		for i := range cfg.Profiles {
			cfg.Profiles[i].Default = false
		}
	}

	cfg.Profiles = append(cfg.Profiles, p)
	return Save(cfg)
}

func RemoveProfile(cfg *Config, name string) error {
	found := false
	filtered := make([]Profile, 0, len(cfg.Profiles))
	for _, p := range cfg.Profiles {
		if p.Name == name {
			found = true
			continue
		}
		filtered = append(filtered, p)
	}
	if !found {
		return fmt.Errorf("profile '%s' not found", name)
	}
	cfg.Profiles = filtered
	return Save(cfg)
}

func SetDefault(cfg *Config, name string) error {
	found := false
	for i := range cfg.Profiles {
		if cfg.Profiles[i].Name == name {
			cfg.Profiles[i].Default = true
			found = true
		} else {
			cfg.Profiles[i].Default = false
		}
	}
	if !found {
		return fmt.Errorf("profile '%s' not found. Run 'auto-message config list' to see available profiles", name)
	}
	return Save(cfg)
}

func MaskToken(token string) string {
	if len(token) <= 4 {
		return "****"
	}
	return "****" + token[len(token)-4:]
}
