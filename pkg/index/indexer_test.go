package index

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseFile(t *testing.T) {
	content := `init() {
    level.callback = true;
    player thread onConnect();
}

onConnect() {
    self waittill("connected");
    player give_weapon("m1911_mp");
    setDvar("g_speed", 190);
}`

	tmpDir := t.TempDir()
	file := filepath.Join(tmpDir, "test.gsc")
	if err := os.WriteFile(file, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	idx := New()
	if err := parseFile(idx, file, "test.gsc"); err != nil {
		t.Fatal(err)
	}

	if _, ok := idx.Functions["init"]; !ok {
		t.Error("init function not found")
	}
	if _, ok := idx.Functions["onConnect"]; !ok {
		t.Error("onConnect function not found")
	}

	if len(idx.Methods["player.give_weapon"]) == 0 {
		t.Error("player.give_weapon method not found")
	}
	if len(idx.Methods["self.waittill"]) == 0 {
		t.Error("self.waittill method not found")
	}
}

func TestGetContext(t *testing.T) {
	lines := strings.Split("line1\nline2\nline3\nline4\nline5", "\n")

	ctx := getContext(lines, 2, 1)
	if len(ctx) != 3 {
		t.Errorf("expected 3 context lines, got %d", len(ctx))
	}
	if ctx[1] != "line3" {
		t.Errorf("expected middle line to be 'line3', got %s", ctx[1])
	}
}

func TestSaveAndLoad(t *testing.T) {
	idx := New()
	idx.Functions["test"] = Function{
		Name: "test",
		File: "test.gsc",
		Line: 1,
	}
	idx.Files = append(idx.Files, "test.gsc")

	tmpFile := filepath.Join(t.TempDir(), "index.json")
	if err := idx.Save(tmpFile); err != nil {
		t.Fatal(err)
	}

	idx2, err := Load(tmpFile)
	if err != nil {
		t.Fatal(err)
	}

	if _, ok := idx2.Functions["test"]; !ok {
		t.Error("test function not found after load")
	}
}
