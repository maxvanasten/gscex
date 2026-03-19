package search

import (
	"fmt"
	"gscex/pkg/index"
	"strings"
)

type Options struct {
	MaxResults   int
	ContextLines int
	FilesOnly    bool
}

type Result struct {
	File    string   `json:"file"`
	Line    int      `json:"line"`
	Content string   `json:"content"`
	Context []string `json:"context"`
}

type Engine struct {
	idx *index.Index
}

func New(idx *index.Index) *Engine {
	return &Engine{idx: idx}
}

func (e *Engine) SearchFunction(name string) (*index.Function, []index.Entry, bool) {
	// Try exact match first (case-insensitive)
	fn, ok := e.idx.Functions[strings.ToLower(name)]
	if !ok {
		fn, ok = e.idx.Functions[name]
	}

	// If no exact match, try fuzzy/substring matching
	if !ok {
		lowerQuery := strings.ToLower(name)
		// Try direct substring match
		for funcName, funcData := range e.idx.Functions {
			if strings.Contains(strings.ToLower(funcName), lowerQuery) {
				fn = funcData
				ok = true
				break
			}
		}

		// If still no match, try without underscores (e.g., "giveweapon" matches "give_weapon")
		if !ok {
			queryNoUnderscore := strings.ReplaceAll(lowerQuery, "_", "")
			for funcName, funcData := range e.idx.Functions {
				funcNameNoUnderscore := strings.ReplaceAll(strings.ToLower(funcName), "_", "")
				if strings.Contains(funcNameNoUnderscore, queryNoUnderscore) {
					fn = funcData
					ok = true
					break
				}
			}
		}
	}

	if !ok {
		return nil, nil, false
	}

	var usages []index.Entry
	for file, raw := range e.idx.Raw {
		lines := strings.Split(raw, "\n")
		for i, line := range lines {
			if strings.Contains(line, fn.Name+"(") {
				trimmed := strings.TrimSpace(line)
				// Skip if it looks like a definition (name() or name() {)
				isDef := strings.HasPrefix(trimmed, fn.Name+"(") &&
					(!strings.Contains(trimmed, ";") || strings.Contains(trimmed, "{"))
				if !isDef {
					context := getContext(lines, i, 3)
					usages = append(usages, index.Entry{
						File:    file,
						Line:    i + 1,
						Content: trimmed,
						Context: context,
					})
				}
			}
		}
	}

	return &fn, usages, true
}

// SearchFunctionsFuzzy returns ALL functions matching the query (for async loading)
func (e *Engine) SearchFunctionsFuzzy(query string, maxResults int) []index.Function {
	lowerQuery := strings.ToLower(query)
	queryNoUnderscore := strings.ReplaceAll(lowerQuery, "_", "")

	var results []index.Function
	seen := make(map[string]bool)

	for funcName, funcData := range e.idx.Functions {
		if seen[funcName] {
			continue
		}

		// Check for match
		match := false

		// Direct substring match
		if strings.Contains(strings.ToLower(funcName), lowerQuery) {
			match = true
		} else if queryNoUnderscore != "" {
			// Underscore-insensitive match
			funcNameNoUnderscore := strings.ReplaceAll(strings.ToLower(funcName), "_", "")
			if strings.Contains(funcNameNoUnderscore, queryNoUnderscore) {
				match = true
			}
		}

		if match {
			results = append(results, funcData)
			seen[funcName] = true

			if maxResults > 0 && len(results) >= maxResults {
				break
			}
		}
	}

	return results
}

func (e *Engine) SearchMethod(entity, method string) []index.Entry {
	var results []index.Entry

	// If method is empty, return all methods for the entity
	if method == "" {
		for key, entries := range e.idx.Methods {
			if strings.HasPrefix(key, entity+".") {
				results = append(results, entries...)
			}
		}
		return results
	}

	// Otherwise, search for specific method
	key := entity + "." + method
	if entries, ok := e.idx.Methods[key]; ok {
		return entries
	}

	// Fuzzy search: find methods containing the search term
	for key, entries := range e.idx.Methods {
		if strings.HasPrefix(key, entity+".") && strings.Contains(key, method) {
			results = append(results, entries...)
		}
	}

	return results
}

func (e *Engine) SearchText(pattern string, opts Options) []Result {
	pattern = strings.ToLower(pattern)
	var results []Result

	for file, raw := range e.idx.Raw {
		lines := strings.Split(raw, "\n")
		for i, line := range lines {
			if strings.Contains(strings.ToLower(line), pattern) {
				if opts.FilesOnly {
					results = append(results, Result{File: file})
					break
				}

				context := getContext(lines, i, opts.ContextLines)
				results = append(results, Result{
					File:    file,
					Line:    i + 1,
					Content: strings.TrimSpace(line),
					Context: context,
				})

				if len(results) >= opts.MaxResults {
					return results
				}
			}
		}
	}

	return results
}

func (e *Engine) ListFiles(pattern string) []string {
	var files []string
	seen := make(map[string]bool)
	pattern = strings.ToLower(pattern)

	for file := range e.idx.Raw {
		if strings.Contains(strings.ToLower(file), pattern) {
			if !seen[file] {
				seen[file] = true
				files = append(files, file)
			}
		}
	}

	return files
}

func (e *Engine) Stats() (files, funcs, methods int) {
	return len(e.idx.Files), len(e.idx.Functions), len(e.idx.Methods)
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

func FormatResult(r Result, highlight string) string {
	content := r.Content
	if highlight != "" {
		content = strings.ReplaceAll(content, highlight, fmt.Sprintf("\033[1;33m%s\033[0m", highlight))
	}
	return fmt.Sprintf("%s:%d\n     %s", r.File, r.Line, content)
}

func FormatFunction(fn *index.Function) string {
	return fmt.Sprintf("%s:%d\n     %s\n     %s",
		fn.File, fn.Line, fn.Signature, strings.Join(fn.Context, "\n     "))
}

func FormatMethod(e index.Entry) string {
	return fmt.Sprintf("%s:%d\n     %s", e.File, e.Line, e.Content)
}
