package git

import (
	"testing"
)

func TestNew(t *testing.T) {
	client := New("https://github.com/test/repo", "main", "/tmp/test")

	if client.repo != "https://github.com/test/repo" {
		t.Errorf("wrong repo: %s", client.repo)
	}
	if client.branch != "main" {
		t.Errorf("wrong branch: %s", client.branch)
	}
	if client.path != "/tmp/test" {
		t.Errorf("wrong path: %s", client.path)
	}
}

func TestExists(t *testing.T) {
	tmpDir := t.TempDir()
	client := New("", "", tmpDir)

	if !client.Exists() {
		t.Error("expected Exists() to be true for existing dir")
	}

	client2 := New("", "", "/nonexistent/path/12345")
	if client2.Exists() {
		t.Error("expected Exists() to be false for non-existing dir")
	}
}

func TestCloneDirExists(t *testing.T) {
	tmpDir := t.TempDir()
	client := New("https://github.com/test/repo", "main", tmpDir)

	err := client.Clone()
	if err == nil {
		t.Error("expected error when directory exists")
	}
}
