package config

import (
	"os"
	"path/filepath"
	"testing"
)

func setupTestConfig(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	orig := os.Getenv("XDG_CONFIG_HOME")
	os.Setenv("XDG_CONFIG_HOME", dir)
	t.Cleanup(func() {
		if orig == "" {
			os.Unsetenv("XDG_CONFIG_HOME")
		} else {
			os.Setenv("XDG_CONFIG_HOME", orig)
		}
	})
	return dir
}

func TestLoad_NoFile(t *testing.T) {
	setupTestConfig(t)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() returned error: %v", err)
	}
	if cfg == nil {
		t.Fatal("Load() returned nil config")
	}
	if len(cfg.Profiles) != 0 {
		t.Fatalf("expected 0 profiles, got %d", len(cfg.Profiles))
	}
}

func TestSaveAndLoad(t *testing.T) {
	setupTestConfig(t)

	cfg := &Config{
		Profiles: []Profile{
			{Name: "test", Token: "123456:ABC-DEF", ChatID: "999", Default: true},
		},
	}

	if err := Save(cfg); err != nil {
		t.Fatalf("Save() returned error: %v", err)
	}

	loaded, err := Load()
	if err != nil {
		t.Fatalf("Load() returned error: %v", err)
	}

	if len(loaded.Profiles) != 1 {
		t.Fatalf("expected 1 profile, got %d", len(loaded.Profiles))
	}
	p := loaded.Profiles[0]
	if p.Name != "test" {
		t.Errorf("expected name 'test', got '%s'", p.Name)
	}
	if p.Token != "123456:ABC-DEF" {
		t.Errorf("expected token '123456:ABC-DEF', got '%s'", p.Token)
	}
	if p.ChatID != "999" {
		t.Errorf("expected chat_id '999', got '%s'", p.ChatID)
	}
	if !p.Default {
		t.Error("expected default to be true")
	}
}

func TestConfigPath(t *testing.T) {
	dir := setupTestConfig(t)

	path, err := ConfigPath()
	if err != nil {
		t.Fatalf("ConfigPath() returned error: %v", err)
	}
	expected := filepath.Join(dir, "auto-message", "config.json")
	if path != expected {
		t.Errorf("expected '%s', got '%s'", expected, path)
	}
}

func TestSave_CreatesDirectory(t *testing.T) {
	setupTestConfig(t)

	cfg := &Config{Profiles: []Profile{{Name: "a", Token: "t", ChatID: "c"}}}
	if err := Save(cfg); err != nil {
		t.Fatalf("Save() returned error: %v", err)
	}

	path, _ := ConfigPath()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Fatal("config file was not created")
	}
}

func TestSave_FilePermissions(t *testing.T) {
	setupTestConfig(t)

	cfg := &Config{Profiles: []Profile{{Name: "a", Token: "secret", ChatID: "c"}}}
	if err := Save(cfg); err != nil {
		t.Fatalf("Save() returned error: %v", err)
	}

	path, _ := ConfigPath()
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("Stat() returned error: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected file permissions 0600, got %o", info.Mode().Perm())
	}
}

func TestAddProfile(t *testing.T) {
	setupTestConfig(t)

	cfg := &Config{}
	p := Profile{Name: "mybot", Token: "123:ABC", ChatID: "42", Default: true}

	if err := AddProfile(cfg, p); err != nil {
		t.Fatalf("AddProfile() returned error: %v", err)
	}

	if len(cfg.Profiles) != 1 {
		t.Fatalf("expected 1 profile, got %d", len(cfg.Profiles))
	}
	if cfg.Profiles[0].Name != "mybot" {
		t.Errorf("expected name 'mybot', got '%s'", cfg.Profiles[0].Name)
	}
}

func TestAddProfile_DuplicateName(t *testing.T) {
	setupTestConfig(t)

	cfg := &Config{}
	p1 := Profile{Name: "mybot", Token: "123:ABC", ChatID: "42"}
	p2 := Profile{Name: "mybot", Token: "456:DEF", ChatID: "99"}

	if err := AddProfile(cfg, p1); err != nil {
		t.Fatalf("first AddProfile() returned error: %v", err)
	}
	if err := AddProfile(cfg, p2); err == nil {
		t.Fatal("expected error for duplicate name, got nil")
	}
}

func TestAddProfile_SetsDefaultClearsPrevious(t *testing.T) {
	setupTestConfig(t)

	cfg := &Config{}
	p1 := Profile{Name: "first", Token: "t1", ChatID: "1", Default: true}
	p2 := Profile{Name: "second", Token: "t2", ChatID: "2", Default: true}

	if err := AddProfile(cfg, p1); err != nil {
		t.Fatalf("first AddProfile() returned error: %v", err)
	}
	if err := AddProfile(cfg, p2); err != nil {
		t.Fatalf("second AddProfile() returned error: %v", err)
	}

	if cfg.Profiles[0].Default {
		t.Error("first profile should not be default after adding second as default")
	}
	if !cfg.Profiles[1].Default {
		t.Error("second profile should be default")
	}
}

func TestRemoveProfile(t *testing.T) {
	setupTestConfig(t)

	cfg := &Config{}
	AddProfile(cfg, Profile{Name: "a", Token: "t1", ChatID: "1"})
	AddProfile(cfg, Profile{Name: "b", Token: "t2", ChatID: "2"})

	if err := RemoveProfile(cfg, "a"); err != nil {
		t.Fatalf("RemoveProfile() returned error: %v", err)
	}

	if len(cfg.Profiles) != 1 {
		t.Fatalf("expected 1 profile, got %d", len(cfg.Profiles))
	}
	if cfg.Profiles[0].Name != "b" {
		t.Errorf("expected remaining profile 'b', got '%s'", cfg.Profiles[0].Name)
	}
}

func TestRemoveProfile_NotFound(t *testing.T) {
	setupTestConfig(t)

	cfg := &Config{}
	AddProfile(cfg, Profile{Name: "a", Token: "t1", ChatID: "1"})

	if err := RemoveProfile(cfg, "nonexistent"); err == nil {
		t.Fatal("expected error for nonexistent profile, got nil")
	}
}

func TestSetDefault(t *testing.T) {
	setupTestConfig(t)

	cfg := &Config{}
	AddProfile(cfg, Profile{Name: "a", Token: "t1", ChatID: "1", Default: true})
	AddProfile(cfg, Profile{Name: "b", Token: "t2", ChatID: "2"})

	if err := SetDefault(cfg, "b"); err != nil {
		t.Fatalf("SetDefault() returned error: %v", err)
	}

	if cfg.Profiles[0].Default {
		t.Error("profile 'a' should not be default")
	}
	if !cfg.Profiles[1].Default {
		t.Error("profile 'b' should be default")
	}
}

func TestSetDefault_NotFound(t *testing.T) {
	setupTestConfig(t)

	cfg := &Config{}
	AddProfile(cfg, Profile{Name: "a", Token: "t1", ChatID: "1"})

	if err := SetDefault(cfg, "nonexistent"); err == nil {
		t.Fatal("expected error for nonexistent profile, got nil")
	}
}

func TestGetDefault_WithDefaultSet(t *testing.T) {
	cfg := &Config{
		Profiles: []Profile{
			{Name: "a", Token: "t1", ChatID: "1"},
			{Name: "b", Token: "t2", ChatID: "2", Default: true},
		},
	}

	p, err := GetDefault(cfg)
	if err != nil {
		t.Fatalf("GetDefault() returned error: %v", err)
	}
	if p.Name != "b" {
		t.Errorf("expected profile 'b', got '%s'", p.Name)
	}
}

func TestGetDefault_FallsToFirst(t *testing.T) {
	cfg := &Config{
		Profiles: []Profile{
			{Name: "a", Token: "t1", ChatID: "1"},
			{Name: "b", Token: "t2", ChatID: "2"},
		},
	}

	p, err := GetDefault(cfg)
	if err != nil {
		t.Fatalf("GetDefault() returned error: %v", err)
	}
	if p.Name != "a" {
		t.Errorf("expected profile 'a', got '%s'", p.Name)
	}
}

func TestGetDefault_Empty(t *testing.T) {
	cfg := &Config{}

	_, err := GetDefault(cfg)
	if err == nil {
		t.Fatal("expected error for empty config, got nil")
	}
}

func TestGetProfile(t *testing.T) {
	cfg := &Config{
		Profiles: []Profile{
			{Name: "a", Token: "t1", ChatID: "1"},
			{Name: "b", Token: "t2", ChatID: "2"},
		},
	}

	p, err := GetProfile(cfg, "b")
	if err != nil {
		t.Fatalf("GetProfile() returned error: %v", err)
	}
	if p.Name != "b" {
		t.Errorf("expected profile 'b', got '%s'", p.Name)
	}
}

func TestGetProfile_NotFound(t *testing.T) {
	cfg := &Config{
		Profiles: []Profile{{Name: "a", Token: "t1", ChatID: "1"}},
	}

	_, err := GetProfile(cfg, "nonexistent")
	if err == nil {
		t.Fatal("expected error for nonexistent profile, got nil")
	}
}

func TestMaskToken(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"123456789:ABCdef", "****Cdef"},
		{"short", "****hort"},
		{"ab", "****"},
		{"", "****"},
		{"1234567890", "****7890"},
	}

	for _, tt := range tests {
		result := MaskToken(tt.input)
		if result != tt.expected {
			t.Errorf("MaskToken(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}
