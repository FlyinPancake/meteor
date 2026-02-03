package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestLoadFilePathConfigExactMatch(t *testing.T) {
	root := t.TempDir()
	project := filepath.Join(root, "project")
	if err := os.MkdirAll(project, 0o755); err != nil {
		t.Fatalf("mkdir project: %v", err)
	}

	configPath := filepath.Join(root, "config.json")
	content, err := json.Marshal(map[string]any{
		project: map[string]any{"showIntro": false},
	})
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if err := os.WriteFile(configPath, content, 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	originalWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	if err := os.Chdir(project); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	defer func() {
		if err := os.Chdir(originalWD); err != nil {
			t.Fatalf("restore wd: %v", err)
		}
	}()

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
	if err := os.MkdirAll(filepath.Join(project, "subdir"), 0o755); err != nil {
		t.Fatalf("mkdir subdir: %v", err)
	}

	configPath := filepath.Join(root, "config.json")
	content, err := json.Marshal(map[string]any{
		project: map[string]any{"showIntro": false},
	})
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if err := os.WriteFile(configPath, content, 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	originalWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	if err := os.Chdir(filepath.Join(project, "subdir")); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	defer func() {
		if err := os.Chdir(originalWD); err != nil {
			t.Fatalf("restore wd: %v", err)
		}
	}()

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
	if err := os.MkdirAll(project, 0o755); err != nil {
		t.Fatalf("mkdir project: %v", err)
	}

	configPath := filepath.Join(root, "config.json")
	content, err := json.Marshal(map[string]any{
		project: map[string]any{"showIntro": false},
	})
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if err := os.WriteFile(configPath, content, 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	originalWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	if err := os.Chdir(root); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	defer func() {
		if err := os.Chdir(originalWD); err != nil {
			t.Fatalf("restore wd: %v", err)
		}
	}()

	c := New()
	if err := c.LoadFile(configPath); err != nil {
		t.Fatalf("LoadFile returned error for no matching path: %v", err)
	}
	if c.ShowIntro != nil {
		t.Fatalf("expected showIntro to remain unset when no match, got %+v", *c.ShowIntro)
	}
}

func TestLoadFilePathConfigGlobMatch(t *testing.T) {
	root := t.TempDir()
	project := filepath.Join(root, "project-a")
	if err := os.MkdirAll(project, 0o755); err != nil {
		t.Fatalf("mkdir project: %v", err)
	}

	configPath := filepath.Join(root, "config.json")
	content, err := json.Marshal(map[string]any{
		filepath.Join(root, "project-*"): map[string]any{"showIntro": false},
	})
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if err := os.WriteFile(configPath, content, 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	originalWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	if err := os.Chdir(project); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	defer func() {
		if err := os.Chdir(originalWD); err != nil {
			t.Fatalf("restore wd: %v", err)
		}
	}()

	c := New()
	if err := c.LoadFile(configPath); err != nil {
		t.Fatalf("LoadFile returned error: %v", err)
	}

	if c.ShowIntro == nil || *c.ShowIntro {
		t.Fatalf("expected showIntro to be false for glob match, got %+v", c.ShowIntro)
	}
}

func TestLoadFilePathConfigGlobMatchReposName(t *testing.T) {
	root := t.TempDir()
	project := filepath.Join(root, "project-a")
	if err := os.MkdirAll(filepath.Join(project, "subdir"), 0o755); err != nil {
		t.Fatalf("mkdir project: %v", err)
	}

	configPath := filepath.Join(root, "repos.json")
	content, err := json.Marshal(map[string]any{
		filepath.Join(root, "project-*"): map[string]any{"showIntro": false},
	})
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if err := os.WriteFile(configPath, content, 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	originalWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	if err := os.Chdir(filepath.Join(project, "subdir")); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	defer func() {
		if err := os.Chdir(originalWD); err != nil {
			t.Fatalf("restore wd: %v", err)
		}
	}()

	c := New()
	if err := c.LoadFile(configPath); err != nil {
		t.Fatalf("LoadFile returned error: %v", err)
	}

	if c.ShowIntro == nil || *c.ShowIntro {
		t.Fatalf("expected showIntro to be false for glob match, got %+v", c.ShowIntro)
	}
}

func TestLoadFilePathConfigGlobMatchSubdir(t *testing.T) {
	root := t.TempDir()
	project := filepath.Join(root, "project-a")
	subdir := filepath.Join(project, "subdir")
	if err := os.MkdirAll(subdir, 0o755); err != nil {
		t.Fatalf("mkdir subdir: %v", err)
	}

	configPath := filepath.Join(root, "config.json")
	content, err := json.Marshal(map[string]any{
		filepath.Join(root, "project-*"): map[string]any{"showIntro": false},
	})
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if err := os.WriteFile(configPath, content, 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	originalWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	if err := os.Chdir(subdir); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	defer func() {
		if err := os.Chdir(originalWD); err != nil {
			t.Fatalf("restore wd: %v", err)
		}
	}()

	c := New()
	if err := c.LoadFile(configPath); err != nil {
		t.Fatalf("LoadFile returned error: %v", err)
	}

	if c.ShowIntro == nil || *c.ShowIntro {
		t.Fatalf("expected showIntro to be false for glob match in subdir, got %+v", c.ShowIntro)
	}
}
