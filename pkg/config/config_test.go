package config

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestLoadFilePathConfigExactMatch(t *testing.T) {
	root := t.TempDir()
	project := filepath.Join(root, "project")
	err := os.MkdirAll(project, 0o755)
	if err != nil {
		t.Fatalf("mkdir project: %v", err)
	}

	configPath := filepath.Join(root, "config.json")
	content := []byte(`{ "` + project + `": { "showIntro": false } }`)
	err = os.WriteFile(configPath, content, 0o644)
	if err != nil {
		t.Fatalf("write config: %v", err)
	}

	originalWD, _ := os.Getwd()
	defer os.Chdir(originalWD)
	_ = os.Chdir(project)

	c := New()
	if err := c.LoadFile(configPath); err != nil {
		t.Fatalf("LoadFile returned error: %v", err)
	}

	if c.ShowIntro == nil || *c.ShowIntro {
		t.Fatalf("expected showIntro to be false, got %+v", c.ShowIntro)
	}
}

func TestLoadFilePathConfigPrefixMatch(t *testing.T) {
	root := t.TempDir()
	project := filepath.Join(root, "project")
	err := os.MkdirAll(filepath.Join(project, "subdir"), 0o755)
	if err != nil {
		t.Fatalf("mkdir subdir: %v", err)
	}

	configPath := filepath.Join(root, "config.json")
	content := []byte(`{ "` + project + `": { "showIntro": false } }`)
	err = os.WriteFile(configPath, content, 0o644)
	if err != nil {
		t.Fatalf("write config: %v", err)
	}

	originalWD, _ := os.Getwd()
	defer os.Chdir(originalWD)
	_ = os.Chdir(filepath.Join(project, "subdir"))

	c := New()
	if err := c.LoadFile(configPath); err != nil {
		t.Fatalf("LoadFile returned error: %v", err)
	}

	if c.ShowIntro == nil || *c.ShowIntro {
		t.Fatalf("expected showIntro to be false for prefix match, got %+v", c.ShowIntro)
	}
}

func TestLoadFilePathConfigNoMatch(t *testing.T) {
	root := t.TempDir()
	project := filepath.Join(root, "project")
	err := os.MkdirAll(project, 0o755)
	if err != nil {
		t.Fatalf("mkdir project: %v", err)
	}

	configPath := filepath.Join(root, "config.json")
	content := []byte(`{ "` + project + `": { "showIntro": false } }`)
	err = os.WriteFile(configPath, content, 0o644)
	if err != nil {
		t.Fatalf("write config: %v", err)
	}

	originalWD, _ := os.Getwd()
	defer os.Chdir(originalWD)
	_ = os.Chdir(root)

	c := New()
	if err := c.LoadFile(configPath); err == nil {
		t.Fatalf("expected error for no matching path")
	} else if !errors.Is(err, ErrNoMatchingPath) {
		t.Fatalf("expected ErrNoMatchingPath, got %v", err)
	}
}
