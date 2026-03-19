package index

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type Function struct {
	Name      string   `json:"name"`
	File      string   `json:"file"`
	Line      int      `json:"line"`
	Signature string   `json:"signature"`
	Context   []string `json:"context"`
}

type Entry struct {
	File    string   `json:"file"`
	Line    int      `json:"line"`
	Content string   `json:"content"`
	Context []string `json:"context"`
}

type Index struct {
	Functions map[string]Function `json:"functions"`
	Methods   map[string][]Entry  `json:"methods"`
	Dvars     map[string][]Entry  `json:"dvars"`
	Files     []string            `json:"files"`
	Raw       map[string]string   `json:"raw"`
}

func New() *Index {
	return &Index{
		Functions: make(map[string]Function),
		Methods:   make(map[string][]Entry),
		Dvars:     make(map[string][]Entry),
		Files:     []string{},
		Raw:       make(map[string]string),
	}
}

var (
	funcDefRegex    = regexp.MustCompile(`^([a-zA-Z_][a-zA-Z0-9_]*)\s*\(([^)]*)\)\s*\{`)
	funcDefRegexAlt = regexp.MustCompile(`^([a-zA-Z_][a-zA-Z0-9_]*)\s*\(([^)]*)\)\s*$`)
	methodCallRegex = regexp.MustCompile(`\b([a-zA-Z_][a-zA-Z0-9_]*)\s+([a-zA-Z_][a-zA-Z0-9_]*)\s*\(`)
	dvarRegex       = regexp.MustCompile(`(?:getDvar|setDvar|getDvarInt|getDvarFloat)\s*\(\s*["']([^"']+)["']`)
)

func Build(root string) (*Index, error) {
	idx := New()

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() || !strings.HasSuffix(path, ".gsc") {
			return nil
		}

		rel, _ := filepath.Rel(root, path)
		idx.Files = append(idx.Files, rel)

		return parseFile(idx, path, rel)
	})

	return idx, err
}

func parseFile(idx *Index, path, rel string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	idx.Raw[rel] = strings.Join(lines, "\n")

	for i, line := range lines {
		// Try standard pattern: funcname() {
		matches := funcDefRegex.FindStringSubmatch(line)
		if matches == nil {
			// Try alternate pattern: funcname() with { on next line
			if i+1 < len(lines) && strings.TrimSpace(lines[i+1]) == "{" {
				matches = funcDefRegexAlt.FindStringSubmatch(line)
			}
		}

		if matches != nil {
			name := matches[1]
			params := matches[2]
			context := getContext(lines, i, 3)

			idx.Functions[name] = Function{
				Name:      name,
				File:      rel,
				Line:      i + 1,
				Signature: name + "(" + params + ")",
				Context:   context,
			}
		}

		if matches := methodCallRegex.FindAllStringSubmatch(line, -1); matches != nil {
			for _, m := range matches {
				entity := m[1]
				method := m[2]
				key := entity + "." + method
				context := getContext(lines, i, 3)

				idx.Methods[key] = append(idx.Methods[key], Entry{
					File:    rel,
					Line:    i + 1,
					Content: strings.TrimSpace(line),
					Context: context,
				})
			}
		}

		if matches := dvarRegex.FindAllStringSubmatch(line, -1); matches != nil {
			for _, m := range matches {
				dvar := m[1]
				context := getContext(lines, i, 3)

				idx.Dvars[dvar] = append(idx.Dvars[dvar], Entry{
					File:    rel,
					Line:    i + 1,
					Content: line,
					Context: context,
				})
			}
		}
	}

	return scanner.Err()
}

func getContext(lines []string, idx, count int) []string {
	start := idx - count
	if start < 0 {
		start = 0
	}
	end := idx + count + 1
	if end > len(lines) {
		end = len(lines)
	}
	return lines[start:end]
}

func (idx *Index) Save(path string) error {
	data, err := json.MarshalIndent(idx, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func Load(path string) (*Index, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var idx Index
	if err := json.Unmarshal(data, &idx); err != nil {
		return nil, err
	}
	return &idx, nil
}
