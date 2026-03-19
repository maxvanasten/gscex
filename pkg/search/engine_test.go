package search

import (
	"gscex/pkg/index"
	"testing"
)

func TestSearchFunction(t *testing.T) {
	idx := index.New()
	idx.Functions["init"] = index.Function{
		Name:      "init",
		File:      "test.gsc",
		Line:      1,
		Signature: "init()",
	}
	idx.Functions["onConnect"] = index.Function{
		Name:      "onConnect",
		File:      "test.gsc",
		Line:      5,
		Signature: "onConnect()",
	}
	idx.Raw["test.gsc"] = "init()\n{\n}\n\nonConnect()\n{\n    init();\n}"

	eng := New(idx)
	fn, usages, ok := eng.SearchFunction("init")

	if !ok {
		t.Fatal("init function not found")
	}
	if fn.Name != "init" {
		t.Errorf("wrong function name: %s", fn.Name)
	}
	if len(usages) == 0 {
		t.Error("expected usage of init()")
	}
}

func TestSearchMethod(t *testing.T) {
	idx := index.New()
	idx.Methods["player.give_weapon"] = []index.Entry{
		{File: "test.gsc", Line: 10, Content: "player give_weapon(\"m1911_mp\");"},
	}

	eng := New(idx)
	results := eng.SearchMethod("player", "give_weapon")

	if len(results) != 1 {
		t.Errorf("expected 1 result, got %d", len(results))
	}
}

func TestSearchText(t *testing.T) {
	idx := index.New()
	idx.Raw["test.gsc"] = "line1\nline2\nweapon give\nline4"
	idx.Raw["test2.gsc"] = "other line\nweapon test"

	eng := New(idx)
	opts := Options{MaxResults: 10, ContextLines: 1, FilesOnly: false}
	results := eng.SearchText("weapon", opts)

	if len(results) != 2 {
		t.Errorf("expected 2 results, got %d", len(results))
	}
}

func TestListFiles(t *testing.T) {
	idx := index.New()
	idx.Raw["maps/mp/test.gsc"] = ""
	idx.Raw["maps/zm/test.gsc"] = ""
	idx.Raw["common.gsc"] = ""

	eng := New(idx)
	files := eng.ListFiles("zm")

	if len(files) != 1 {
		t.Errorf("expected 1 zm file, got %d", len(files))
	}
}

func TestStats(t *testing.T) {
	idx := index.New()
	idx.Files = append(idx.Files, "test1.gsc", "test2.gsc")
	idx.Functions["f1"] = index.Function{}
	idx.Functions["f2"] = index.Function{}
	idx.Methods["m1"] = []index.Entry{{}}

	eng := New(idx)
	files, funcs, methods := eng.Stats()

	if files != 2 {
		t.Errorf("expected 2 files, got %d", files)
	}
	if funcs != 2 {
		t.Errorf("expected 2 functions, got %d", funcs)
	}
	if methods != 1 {
		t.Errorf("expected 1 method, got %d", methods)
	}
}
